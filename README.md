レイヤードアーキテクチャにおけるトランザクションの実装案

### [ContextにTxオブジェクトを詰めるパターン](./context-pattern)

#### 概要
- ContextにTxオブジェクトを詰める
- RepositoryでContextのValueを参照し、TxオブジェクトがあればTxオブジェクトを、ない場合はDIされた素のDBオブジェクトを利用する。

(usecase単位でトランザクション処理が必要な部分だけラップするか、middlewareで各エンドポイント全体をラップするかは選択)

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

// context-pattern/infra/mysql/tx_manager.go
func (t *txManager) Transaction(ctx context.Context, f func(context.Context) error) error {
    tx, err := t.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer func() {
        // (recovery process...)
        // success
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

#### Pros/Cons
- Pros
  - 同一I/FのRepositoryで、Transactionで実行するかどうか切り分けられる
  - contextさえ受け渡していればどこでもTxオブジェクトを取り出せる
- Cons
  - contextの乱用感が否めない
  - どこでDBアクセスが発生するのか/Transactionが使用されているのか分かりにくい
  - ReadWriteTransactionとReadOnlyTransactionを使い分ける実装が難しく、見た目上も分かりにくい
