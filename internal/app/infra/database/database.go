package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func NewPostgresDatabase() *sql.DB {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)

		pg, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		pg.SetMaxOpenConns(25)
		pg.SetMaxIdleConns(25)

		if err := pg.Ping(); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		pg.Exec(
			peopledb.InsertPersonQuery,
			uuid.NewString(),
			"leorcvargas",
			"Leonardo Vargas",
			"1970-01-01",
			"Go, Node.js",
			"leorcvargas "+"Leonardo Vargas"+"Go, Node.js",
		)

		db = pg
	})

	return db
}
