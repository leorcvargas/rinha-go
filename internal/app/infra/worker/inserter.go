package worker

import (
	"arena"
	"context"
	"database/sql"
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
			i.cache.Set(person.ID, &person)

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
	stmt, err := i.db.Prepare(peopledb.InsertPersonQuery)
	if err != nil {
		log.Errorf("Error preparing batch: %v", err)
		return err
	}
	defer stmt.Close()

	tx, err := i.db.Begin()
	if err != nil {
		log.Errorf("Error starting transaction: %v", err)
		return err
	}

	for index := 0; index < batchLength; index++ {
		person := batch[index]

		if person.ID == "" {
			continue
		}

		_, err := tx.Stmt(stmt).Exec(
			person.ID,
			person.Nickname,
			person.Name,
			person.Birthdate,
			strings.Join(person.Stack, ","),
			person.Nickname+person.Name+strings.Join(person.Stack, ""),
		)

		if err != nil {
			log.Errorf("Error inserting batch: %v", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Error committing transaction: %v", err)
		return err
	}

	return nil

	// valueStrings := make([]string, batchLength, batchLength)
	// valueArgs := make([]interface{}, batchLength*6, batchLength*6)

	// for index := 0; index < batchLength; index++ {
	// 	person := batch[index]

	// 	if person.ID == "" {
	// 		continue
	// 	}

	// 	valueStrings[index] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", index*6+1, index*6+2, index*6+3, index*6+4, index*6+5, index*6+6)
	// 	valueArgs[index*6] = person.ID
	// 	valueArgs[index*6+1] = person.Nickname
	// 	valueArgs[index*6+2] = person.Name
	// 	valueArgs[index*6+3] = person.Birthdate
	// 	valueArgs[index*6+4] = person.StackString()
	// 	valueArgs[index*6+5] = person.Nickname + person.Name + strings.Join(person.Stack, "")
	// }

	// stmt := "INSERT INTO people (id, nickname, name, birthdate, stack, search) VALUES "
	// for i := 0; i < len(valueStrings); i++ {
	// 	if i == 0 {
	// 		stmt += valueStrings[i]
	// 	} else {
	// 		stmt += "," + valueStrings[i]
	// 	}
	// }

	// _, err := i.db.Exec(stmt, valueArgs...)
	// if err != nil {
	// 	log.Errorf("Error inserting batch: %v", err)
	// 	return err
	// }

	// return nil
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
