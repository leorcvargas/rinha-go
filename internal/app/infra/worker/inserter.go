package worker

import (
	"arena"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/log"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/pubsub"
)

type Inserter struct {
	insertChan chan people.Person
	db         *sql.DB
	cache      *peopledb.PeopleDbCache
}

const maxBatchSize = 10000

func (i *Inserter) Run() {
	a := arena.NewArena()

	batch := arena.MakeSlice[people.Person](a, maxBatchSize, maxBatchSize)
	batchLen := 0

	tickProcess := time.Tick(2 * time.Second)
	tickClear := time.Tick(1 * time.Minute)

	for {
		select {
		case person := <-i.insertChan:
			batch[batchLen] = person
			batchLen++

			// if batchLen >= maxBatchSize {
			// 	i.processBatch(batch, batchLen)
			// 	a.Free()
			// 	a = arena.NewArena()
			// 	batch = arena.MakeSlice[people.Person](a, maxBatchSize, maxBatchSize)
			// 	batchLen = 0
			// }

		case <-tickProcess:
			if batchLen > 0 {
				i.processBatch(batch, batchLen)
				batch = arena.MakeSlice[people.Person](a, maxBatchSize, maxBatchSize)
				batchLen = 0
			}

		case <-tickClear:
			log.Info("Clear tick...")
			if batchLen > 0 {
				i.processBatch(batch, batchLen)
			}

			a.Free()
			a = arena.NewArena()
			batch = arena.MakeSlice[people.Person](a, maxBatchSize, maxBatchSize)
			batchLen = 0
		}

	}
}

func (i *Inserter) processBatch(batch []people.Person, batchLength int) error {
	err := i.insertBatch(batch, batchLength)
	if err != nil {
		return err
	}

	payload, err := sonic.MarshalString(batch[:batchLength])
	if err != nil {
		log.Errorf("Error marshalling batch: %v", err)
		return err
	}

	err = i.cache.
		Cache().
		Do(
			context.Background(),
			i.cache.
				Cache().
				B().
				Publish().
				Channel(pubsub.EventPersonInsert).
				Message(payload).
				Build(),
		).
		Error()
	if err != nil {
		log.Errorf("Error publishing batch: %v", err)
		return err
	}

	return nil
}

func (i *Inserter) insertBatch(batch []people.Person, batchLength int) error {
	totalCols := 5

	valueStrings := make([]string, batchLength, batchLength)
	valueArgs := make([]interface{}, batchLength*totalCols, batchLength*totalCols)

	for index := 0; index < batchLength; index++ {
		person := batch[index]

		if person.ID == "" {
			continue
		}

		colIndex := index * totalCols

		valueStrings[index] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", colIndex+1, colIndex+2, colIndex+3, colIndex+4, colIndex+5)
		valueArgs[colIndex] = person.ID
		valueArgs[colIndex+1] = person.Nickname
		valueArgs[colIndex+2] = person.Name
		valueArgs[colIndex+3] = person.Birthdate
		valueArgs[colIndex+4] = person.StackString()
	}

	stmt := "INSERT INTO people (id, nickname, name, birthdate, stack) VALUES "
	for i := 0; i < len(valueStrings); i++ {
		if i == 0 {
			stmt += valueStrings[i]
		} else {
			stmt += "," + valueStrings[i]
		}
	}

	_, err := i.db.Exec(stmt, valueArgs...)
	if err != nil {
		log.Errorf("Error inserting batch: %v", err)
		return err
	}

	return nil
}

func NewInserter(
	insertChan chan people.Person,
	db *sql.DB,
	cache *peopledb.PeopleDbCache,
) *Inserter {
	return &Inserter{
		insertChan: insertChan,
		db:         db,
		cache:      cache,
	}
}
