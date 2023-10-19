package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Addr     string
	User     string
	Password string
	DB       string
}

func NewDB(c *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", c.User, c.Password, c.Addr, c.DB)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
