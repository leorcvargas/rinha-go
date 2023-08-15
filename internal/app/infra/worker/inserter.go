package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/pubsub"
)

const batchSize = 20

type Inserter struct {
	insertChan chan people.Person
	db         *sql.DB
	cache      *peopledb.PeopleDbCache
}

func (i *Inserter) Run() {
	batch := i.makeEmptyBatch()

	for {
		select {
		case person := <-i.insertChan:
			batch = append(batch, person)

		case <-time.Tick(5 * time.Second):
			i.processBatch(batch)
			batch = i.makeEmptyBatch()
		}
	}
}

func (*Inserter) makeEmptyBatch() []people.Person {
	return make([]people.Person, 0, batchSize)
}

func (i *Inserter) processBatch(batch []people.Person) error {
	if len(batch) == 0 {
		return nil
	}

	err := i.insertBatch(batch)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(batch)
	if err != nil {
		log.Printf("Error marshalling batch: %v", err)
		return err
	}

	i.cache.Cache().Publish(
		context.Background(),
		pubsub.EventPersonInsert,
		payload,
	)

	return nil
}

func (i *Inserter) insertBatch(batch []people.Person) error {
	// memory := arena.NewArena()
	// defer memory.Free()

	// valueStrings := arena.MakeSlice[string](memory, 0, len(batch))
	// valueArgs := arena.MakeSlice[interface{}](memory, 0, len(batch)*5)

	valueStrings := make([]string, 0, len(batch))
	valueArgs := make([]interface{}, 0, len(batch)*5)

	for i, person := range batch {
		if person.ID == uuid.Nil {
			continue
		}

		valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		valueArgs[i*5+1] = person.ID
		valueArgs[i*5+2] = person.Nickname
		valueArgs[i*5+3] = person.Name
		valueArgs[i*5+4] = person.Birthdate
		valueArgs[i*5+5] = strings.Join(person.Stack, ",")

		// valueStrings = append(
		// 	valueStrings,
		// 	fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5),
		// )
		// valueArgs = append(
		// 	valueArgs,
		// 	person.ID,
		// 	person.Nickname,
		// 	person.Name,
		// 	person.Birthdate,
		// 	strings.Join(person.Stack, ","),
		// )
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
		log.Printf("Error inserting batch: %v", err)
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
