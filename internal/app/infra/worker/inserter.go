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

type Inserter struct {
	insertChan chan people.Person
	db         *sql.DB
	cache      *peopledb.PeopleDbCache
	batch      []people.Person
}

func (i *Inserter) Run() {
	for {
		select {
		case person := <-i.insertChan:
			i.batch = append(i.batch, person)

		case <-time.Tick(5 * time.Second):
			if len(i.batch) == 0 {
				continue
			}

			i.processBatch()
			i.clearBatch()
		}
	}
}

func (i *Inserter) clearBatch() {
	i.batch = make([]people.Person, 0)
}

func (i *Inserter) processBatch() error {
	err := i.insertBatch()
	if err != nil {
		return err
	}

	payload, err := json.Marshal(i.batch)
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

func (i *Inserter) insertBatch() error {
	batchLength := len(i.batch)

	valueStrings := make([]string, batchLength, batchLength)
	valueArgs := make([]interface{}, batchLength*5, batchLength*5)

	for i, person := range i.batch {
		if person.ID == uuid.Nil {
			continue
		}

		valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		valueArgs[i*5] = person.ID
		valueArgs[i*5+1] = person.Nickname
		valueArgs[i*5+2] = person.Name
		valueArgs[i*5+3] = person.Birthdate
		valueArgs[i*5+4] = strings.Join(person.Stack, ",")
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
		batch:      make([]people.Person, 0),
	}
}
