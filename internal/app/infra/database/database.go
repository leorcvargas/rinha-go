package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/config"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

var (
	db   *pgxpool.Pool
	once sync.Once
)

func NewPostgresDatabase(config *config.Config) *pgxpool.Pool {
	once.Do(func() {
		connUrl := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)

		poolConfig, err := pgxpool.ParseConfig(connUrl)
		if err != nil {
			log.Fatalln("Unable to parse connection url:", err)
		}

		poolConfig.MinConns = 5

		db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			log.Fatalln("Unable to create connection pool:", err)
		}

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
				person.StackStr(),
				person.SearchStr(),
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
