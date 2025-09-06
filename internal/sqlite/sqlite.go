package sqlite

import (
	"context"
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var once sync.Once

func InitDB(ctx context.Context, dsn string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", dsn)
		if err != nil {
			panic(err)
		}
		db.SetMaxOpenConns(1)
	})

	return db.PingContext(ctx)
}

func GetDB() *sql.DB {
	return db
}
