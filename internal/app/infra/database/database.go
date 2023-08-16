package database

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func NewPostgresDatabase() *sql.DB {
	once.Do(func() {
		dsn := "host=db user=postgres password=postgres dbname=rinha port=5432 sslmode=disable"

		pg, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		pg.SetMaxOpenConns(25)
		pg.SetMaxIdleConns(25)

		if err := pg.Ping(); err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		db = pg
	})

	return db
}
