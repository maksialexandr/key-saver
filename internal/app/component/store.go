package component

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"person-key-saver/internal/app/application/configure"
	"time"
)

type Store struct {
	Db *sql.DB
}

func NewStore() *Store {
	return &Store{}
}

func (this *Store) Connect(config *configure.Config) error {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Store.User,
		config.Store.Pwd,
		config.Store.Dsn,
		config.Store.Port,
		config.Store.DataBase,
	))

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	this.Db = db

	return nil
}
