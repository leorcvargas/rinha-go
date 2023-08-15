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

const batchSize = 10

type Inserter struct {
	insertChan chan people.Person
	db         *sql.DB
	cache      *peopledb.PeopleDbCache
}

func (i *Inserter) Run() {
	batchSize := 5

	batch := i.makeEmptyBatch()

	for {
		select {
		case person := <-i.insertChan:
			batch = append(batch, person)

			if len(batch) == batchSize {
				i.processBatch(batch)
				batch = i.makeEmptyBatch()
			}

		case <-time.Tick(10 * time.Second):
			if len(batch) > 0 {
				i.processBatch(batch)
				batch = i.makeEmptyBatch()
			}
		}
	}
}

func (*Inserter) makeEmptyBatch() []people.Person {
	return make([]people.Person, 0, batchSize)
}

func (i *Inserter) processBatch(batch []people.Person) {
	i.insertBatch(batch)

	payload, err := json.Marshal(batch)
	if err != nil {
		log.Printf("Error marshalling batch: %v", err)
		return
	}

	i.cache.Cache().Publish(
		context.Background(),
		pubsub.EventPersonInsert,
		payload,
	)
}

func (i *Inserter) insertBatch(batch []people.Person) {
	valueStrings := make([]string, 0, len(batch))
	valueArgs := make([]interface{}, 0, len(batch)*5)

	for i, person := range batch {
		if person.ID == uuid.Nil {
			continue
		}

		valueStrings = append(
			valueStrings,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5),
		)
		valueArgs = append(
			valueArgs,
			person.ID,
			person.Nickname,
			person.Name,
			person.Birthdate,
			strings.Join(person.Stack, ","),
		)
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
	}
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
