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
	batchCap := batchSize * 1024
	batch := arena.MakeSlice[people.Person](a, batchSize, batchCap)
	at := 0

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
		case <-time.Tick(5 * time.Second):
			if at > 0 {
				insertBatch(batch, db)

				a.Free()
				at = 0
				a = arena.NewArena()
				batch = arena.MakeSlice[people.Person](a, batchSize, batchSize)
			}
		}
	}
}

func insertBatch(batch []people.Person, db *sql.DB) {
	bulkInsert := "INSERT INTO people (id, nickname, name, birthdate, stack) VALUES "

	for i, person := range batch {
		strStack := strings.Join(person.Stack, ",")
		bulkInsert += "(" +
			person.ID.String() + ", " +
			person.Nickname + ", " +
			person.Name + ", " +
			person.Birthdate + ", " +
			strStack + ")"

		if i != len(batch)-1 {
			bulkInsert += ", "
		}
	}

	_, err := db.Exec(bulkInsert)
	if err != nil {
		log.Printf("Error inserting batch: %v", err)
	}

	// tx, err := db.Begin()
	// if err != nil {
	// 	panic(err)
	// }

	// stmt, err := tx.Prepare(InsertPersonQuery)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, person := range batch {
	// 	strStack := strings.Join(person.Stack, ",")

	// 	_, err := stmt.Exec(
	// 		person.ID,
	// 		person.Nickname,
	// 		person.Name,
	// 		person.Birthdate,
	// 		strStack,
	// 	)
	// 	if err != nil {
	// 		log.Printf("Error inserting person: %v", err)
	// 	}
	// }

	// err = tx.Commit()
	// if err != nil {
	// 	log.Printf("Error committing insert transaction: %v", err)
	// }
}
