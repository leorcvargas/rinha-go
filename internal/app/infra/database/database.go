package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

var (
	db   *pgxpool.Pool
	once sync.Once
)

func NewPostgresDatabase() *pgxpool.Pool {
	once.Do(func() {
		connUrl := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		poolConfig, err := pgxpool.ParseConfig(connUrl)
		if err != nil {
			log.Fatalln("Unable to parse connection url:", err)
		}

		db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatalln("Unable to create connection pool:", err)
		}

		// dsn := fmt.Sprintf(
		// 	"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		// 	os.Getenv("DB_HOST"),
		// 	os.Getenv("DB_USER"),
		// 	os.Getenv("DB_PASSWORD"),
		// 	os.Getenv("DB_NAME"),
		// 	os.Getenv("DB_PORT"),
		// )

		// pg, err := sql.Open("postgres", dsn)
		// if err != nil {
		// 	log.Fatalf("Failed to connect to database: %v", err)
		// }

		// pg.SetMaxOpenConns(25)
		// pg.SetMaxIdleConns(25)

		// warmup
		var ids []string
		for i := 0; i < 10; i++ {
			person := people.NewPerson(
				fmt.Sprintf("nickname-%d", i),
				fmt.Sprintf("name-%d", i),
				"1970-01-01",
				[]string{"tag1", "tag2"},
			)
			ids = append(ids, person.ID)
			_, err := db.Exec(
				context.Background(),
				peopledb.InsertPersonQuery,
				person.ID,
				person.ID[:30],
				person.Name,
				person.Birthdate,
				person.StackString(),
				"somestring",
			)
			if err != nil {
				log.Fatalf("Failed to warmup database: %v", err)
			}
		}

		for _, id := range ids {
			_, err := db.Exec(
				context.Background(),
				"DELETE FROM people WHERE id = $1",
				id,
			)
			if err != nil {
				log.Fatalf("Failed to delete warmup data from the database: %v", err)
			}
		}

		if err := db.Ping(context.Background()); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	})

	return db
}
