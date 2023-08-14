package peopledb

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

const batchTimerAmount = 5 * time.Second

func Worker(insertChan chan people.Person, db *sql.DB) {
	batchSize := 20
	batch := make([]people.Person, 0, batchSize)

	batchInsertTimer := time.NewTimer(batchTimerAmount)

	for {
		select {
		case person := <-insertChan:
			batch = append(batch, person)

			if len(batch) == batchSize {
				insertBatch(batch, db)
				batch = make([]people.Person, 0, batchSize)
			}
		case <-batchInsertTimer.C:
			if len(batch) > 0 {
				insertBatch(batch, db)
				batch = make([]people.Person, 0, batchSize)
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
