package worker

import (
	"arena"
	"context"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

type Inserter struct {
	insertChan chan people.Person
	db         *pgxpool.Pool
	cache      *peopledb.PeopleDbCache
}

const maxBatchSize = 10000

func (i *Inserter) Run() {
	a := arena.NewArena()

	batch := arena.MakeSlice[people.Person](a, maxBatchSize, maxBatchSize)
	batchLen := 0

	tickProcess := time.Tick(10 * time.Second)
	tickClear := time.Tick(3 * time.Minute)

	for {
		select {
		case person := <-i.insertChan:
			batch[batchLen] = person
			batchLen++

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

	// payload, err := sonic.MarshalString(batch[:batchLength])
	// if err != nil {
	// 	log.Errorf("Error marshalling batch: %v", err)
	// 	return err
	// }

	// err = i.cache.
	// 	Cache().
	// 	Do(
	// 		context.Background(),
	// 		i.cache.
	// 			Cache().
	// 			B().
	// 			Publish().
	// 			Channel(pubsub.EventPersonInsert).
	// 			Message(payload).
	// 			Build(),
	// 	).
	// 	Error()
	// if err != nil {
	// 	log.Errorf("Error publishing batch: %v", err)
	// 	return err
	// }

	return nil
}

func (i *Inserter) insertBatch(batch []people.Person, batchLength int) error {
	dbBatch := &pgx.Batch{}

	query := "INSERT INTO people (id, nickname, name, birthdate, stack) VALUES ($1, $2, $3, $4, $5)"
	for index := 0; index < batchLength; index++ {
		person := batch[index]

		if person.ID == "" {
			continue
		}

		dbBatch.Queue(
			query,
			person.ID,
			person.Nickname,
			person.Name,
			person.Birthdate,
			person.StackString(),
			person.Name+person.Nickname+person.StackString(),
		)
	}

	batchResults := i.db.SendBatch(context.Background(), dbBatch)
	defer batchResults.Close()

	for index := 0; index < batchLength; index++ {
		_, err := batchResults.Exec()
		if err != nil {
			log.Errorf("Error inserting batch: %v", err)
			return err
		}
	}

	return nil
}

func NewInserter(
	insertChan chan people.Person,
	db *pgxpool.Pool,
	cache *peopledb.PeopleDbCache,
) *Inserter {
	return &Inserter{
		insertChan: insertChan,
		db:         db,
		cache:      cache,
	}
}
