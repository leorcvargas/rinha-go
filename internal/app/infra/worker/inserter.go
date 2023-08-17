package worker

import (
	"arena"
	"context"
	"database/sql"
	"fmt"
	"strings"
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

const maxBatchSize = 100

func (i *Inserter) Run() {
	a := arena.NewArena()

	batch := arena.MakeSlice[people.Person](a, maxBatchSize, maxBatchSize)
	batchLen := 0

	tickProcess := time.Tick(5 * time.Second)
	tickClear := time.Tick(4 * time.Minute)

	for {
		select {
		case person := <-i.insertChan:
			batch[batchLen] = person
			batchLen++

			if batchLen >= maxBatchSize {
				i.processBatch(batch, batchLen)
				a.Free()
				a = arena.NewArena()
				batch = arena.MakeSlice[people.Person](a, maxBatchSize, maxBatchSize)
				batchLen = 0
			}

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

	payload, err := sonic.Marshal(batch[:batchLength])
	if err != nil {
		log.Errorf("Error marshalling batch: %v", err)
		return err
	}

	i.cache.Cache().Publish(
		context.Background(),
		pubsub.EventPersonInsert,
		payload,
	)

	return nil
}

func (i *Inserter) insertBatch(batch []people.Person, batchLength int) error {
	valueStrings := make([]string, batchLength, batchLength)
	valueArgs := make([]interface{}, batchLength*5, batchLength*5)

	for index := 0; index < batchLength; index++ {
		person := batch[index]

		if person.ID == "" {
			continue
		}

		valueStrings[index] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", index*5+1, index*5+2, index*5+3, index*5+4, index*5+5)
		valueArgs[index*5] = person.ID
		valueArgs[index*5+1] = person.Nickname
		valueArgs[index*5+2] = person.Name
		valueArgs[index*5+3] = person.Birthdate
		valueArgs[index*5+4] = strings.Join(person.Stack, ",")
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
