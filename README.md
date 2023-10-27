レイヤードアーキテクチャにおけるトランザクションの実装案

### [ContextにTxオブジェクトを詰めるパターン](./context-pattern)

<details><summary>details</summary>

#### 概要
- **ContextにTxオブジェクトを詰める**パターン
- **RepositoryでContextのValueを参照**し、TxオブジェクトがあればTxオブジェクトを、ない場合はDIされた素のDBオブジェクトを利用する。
(usecase単位でトランザクション処理が必要な部分だけラップするか、middlewareで各エンドポイント全体をラップするかは選択)

#### Pros/Cons
- Pros
  - RepositoryでContextさえ受け取っておけば、トランザクション内で実行するかどうか外部から指定できる
- Cons
  - トランザクション内の処理なのかシグニチャで判別できない
  - ReadOnly/ReadWriteなトランザクションを使い分けるのが少し実装大変

#### 実装

```go
// context-pattern/usecase/user.go
func (i *userInteractor) UpdateName(ctx context.Context, userID, name string) error {
    if err := i.txManager.Transaction(ctx, func(ctx context.Context) error {
        // ...
    }); err != nil {
        return err
    }
    return nil
}

// context-pattern/domain/transaction/tx_manager.go
type TxManager interface {
    Transaction(ctx context.Context, f func(context.Context) error) error
}

// context-pattern/infra/mysql/tx_manager.go
func (t *txManager) Transaction(ctx context.Context, f func(context.Context) error) error {
    tx, err := t.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        // (recovery process...)
        if e := tx.Commit(); e != nil {
            slog.ErrorContext(ctx, "failed to MySQL Commit")
        }
    }()

    // ContextにTxをセット
    ctx = xcontext.WithValue[xcontext.MysqlTx, *sqlx.Tx](ctx, tx)
    err = f(ctx)
    if err != nil {
        return err
    }
    return nil
}

// context-pattern/infra/repository/user.go
func (r *userRepository) getMysqlDB(ctx context.Context) infra.MysqlDB {
    // contextにtxオブジェクトが存在すればそれを返却する
    if tx, ok := xcontext.Value[xcontext.MysqlTx, *sqlx.Tx](ctx); ok {
        return tx
    }
    // contextにtxオブジェクトが存在しなければDIされたdbを返却する
    return r.db
}
```

```shell
$ docker compose up -d
$ run-context-pattern
```
</details>

### [Txオブジェクトを抽象化しusecase層にDIするパターン](./di-pattern)



<details><summary>details</summary>

#### 概要
- **Txオブジェクトを抽象化**し、usecase層で扱えるように**DIで注入する**パターン
- **ReadOnlyとReadWriteでTxオブジェクトの抽象を分ける**

#### Pros/Cons
- Pros
  - ReadOnlyかReadWriteかをusecase層で扱えることで、**効率的なTransaction**の貼り方を行える
  - 関数のシグニチャを見ただけで、その処理がどのようなトランザクション内で実行されることを期待しているのかが分かる
  - Repositoryの引数にTxオブジェクトを受け取るように設定できることで、**トランザクションの開始漏れがなくなる**
- Cons
  - トランザクション内/外で実行するRepositoryのシグニチャが異なる（Txオブジェクトを受け取るかどうか）ので、Repositoryの実装が複雑になる可能性がある (プロジェクト内でRepository呼び出しは必ずトランザクション内で行うという合意が取れていればそこまでデメリットにならない気がしている)

#### 実装

```go
// di-pattern/usecase/user.go
func (i *userInteractor) GetUser(ctx context.Context, userID string) (*entity.User, error) {
    var user *entity.User 
    i.txManager.ReadOnlyTransaction(ctx, func(ctx context.Context, tx transaction.ROTx) error {
        // ...
    })
    return user, nil
}

func (i *userInteractor) UpdateName(ctx context.Context, userID, name string) error {
    i.txManager.ReadWriteTransaction(ctx, func(ctx context.Context, tx transaction.RWTx) error {
        // ...
    })
    return nil
}

// di-pattern/domain/transaction/tx_manager.go
type ROTx interface {
    ROTxImpl()
}

type RWTx interface {
    ROTx
    RWTxImpl()
}

type TxManager interface {
    ReadOnlyTransaction(ctx context.Context, f func(ctx context.Context, tx ROTx) error) error
    ReadWriteTransaction(ctx context.Context, f func(ctx context.Context, tx RWTx) error) error
}

// di-pattern/infra/mysql/tx.go
type ROTx interface {
    GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type RWTx interface {
    ROTx
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type rwTx struct {
    *sqlx.Tx
}

func (tx *rwTx) ROTxImpl() {}
func (tx *rwTx) RWTxImpl() {}

func ExtractRWTx(_tx transaction.RWTx) (RWTx, error) {
    tx, ok := _tx.(*rwTx)
    if !ok {
        return nil, errors.New("mysql RWTx is invalid")
    }
    return tx, nil
}

type roTx struct {
    // MysqlにはReadOnlyなTxオブジェクトが存在しない
    *sqlx.Tx
}

func (tx *roTx) ROTxImpl() {}

func ExtractROTx(_tx transaction.ROTx) (ROTx, error) {
    switch tx := _tx.(type) {
    case *roTx:
        return tx, nil
    case *rwTx: // ReadWriteTransaction内での呼び出しも許可する
        return tx, nil
    }
    return nil, errors.New("mysql ROTx is invalid")
}

// di-pattern/infra/mysql/tx_manager.go
func (t *txManager) ReadWriteTransaction(ctx context.Context, f func(context.Context, transaction.RWTx) error) error {
    tx, err := t.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        // (recovery process...)
        if e := tx.Commit(); e != nil {
            slog.ErrorContext(ctx, "failed to MySQL Commit")
        }
    }()

    // ReadWriteTransactionを関数に渡す
    err = f(ctx, &rwTx{tx})
    if err != nil {
        return err
    }
    return nil
    }

func (t *txManager) ReadOnlyTransaction(ctx context.Context, f func(context.Context, transaction.ROTx) error) error {
    tx, err := t.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        // (recovery process...)
        if e := tx.Commit(); e != nil {
            slog.ErrorContext(ctx, "failed to MySQL Commit")
        }
	}()

    // ReadOnlyTransactionを関数に渡す
    err = f(ctx, &roTx{tx})
    if err != nil {
        return err
    }
    return nil
}

// di-pattern/infra/repository/user.go
func (r *userRepository) SelectByPK(ctx context.Context, _tx transaction.ROTx, userID string) (*entity.User, error) {
    tx, err := mysql.ExtractROTx(_tx)
    if err != nil {
        return nil, err
    }

    var user User
    if err := tx.GetContext(ctx, &user, "SELECT * FROM users WHERE user_id = ?", userID); err != nil {
        return nil, err
    }
    return user.toEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, _tx transaction.RWTx, e *entity.User) error {
    tx, err := mysql.ExtractRWTx(_tx)
	if err != nil {
        return err
    }

    if _, err := tx.ExecContext(ctx, "UPDATE users SET name = ? WHERE user_id = ?", e.Name, e.UserID); err != nil {
        return err
    }
    return nil
}
```

```shell
$ docker compose up -d
$ run-di-pattern
```
</details>

個人的には後者がかっこいい気がしている
