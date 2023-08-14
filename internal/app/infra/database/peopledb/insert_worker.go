package peopledb

import (
	"arena"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

const batchTimerAmount = 10 * time.Second

func Worker(insertChan chan people.Person, db *sql.DB) {
	batchSize := 20

	a := arena.NewArena()
	batch := arena.MakeSlice[people.Person](a, batchSize, batchSize)
	at := 0

	batchInsertTimer := time.NewTimer(batchTimerAmount)

	for {
		select {
		case person := <-insertChan:
			batch[at] = person
			at += 1

			if at == batchSize {
				insertBatch(batch, db)

				a.Free()
				at = 0
				a = arena.NewArena()
				batch = arena.MakeSlice[people.Person](a, batchSize, batchSize)
			}
		case <-batchInsertTimer.C:
			if at > 0 {
				insertBatch(batch, db)

				a.Free()
				at = 0
				a = arena.NewArena()
				batch = arena.MakeSlice[people.Person](a, batchSize, batchSize)
			}
			batchInsertTimer.Reset(batchTimerAmount)
		}
	}
}

func insertBatch(batch []people.Person, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	stmt, err := tx.Prepare(InsertPersonQuery)
	if err != nil {
		panic(err)
	}

	for _, person := range batch {
		strStack := strings.Join(person.Stack, ",")

		_, err := stmt.Exec(
			person.ID,
			person.Nickname,
			person.Name,
			person.Birthdate,
			strStack,
		)
		if err != nil {
			log.Printf("Error inserting person: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing insert transaction: %v", err)
	}
}
