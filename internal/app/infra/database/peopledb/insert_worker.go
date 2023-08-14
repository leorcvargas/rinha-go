package peopledb

import (
	"arena"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
	valueStrings := make([]string, 0, len(batch))
	valueArgs := make([]interface{}, 0, len(batch)*5)

	for i, person := range batch {
		if person.ID == uuid.Nil {
			continue
		}

		valuesStrCh := make(chan []string)
		valuesArgCh := make(chan []interface{})

		go func() {
			defer close(valuesStrCh)
			valuesStrCh <- append(
				valueStrings,
				fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5),
			)
		}()

		go func() {
			defer close(valuesArgCh)
			valuesArgCh <- append(
				valueArgs,
				person.ID,
				person.Nickname,
				person.Name,
				person.Birthdate,
				strings.Join(person.Stack, ","),
			)
		}()

		valueStrings = <-valuesStrCh
		valueArgs = <-valuesArgCh
	}

	stmt := "INSERT INTO people (id, nickname, name, birthdate, stack) VALUES "
	for i := 0; i < len(valueStrings); i++ {
		if i == 0 {
			stmt += valueStrings[i]
		} else {
			stmt += "," + valueStrings[i]
		}
	}

	_, err := db.Exec(stmt, valueArgs...)
	if err != nil {
		log.Printf("Error inserting batch: %v", err)
	}
}
