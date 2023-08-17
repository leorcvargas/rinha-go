package peopledb

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hashicorp/go-memdb"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

var ErrMemDbNotFound = errors.New("memdb not found")

type MemDb struct {
	db *memdb.MemDB
}

type PersonMem struct {
	Key string
	people.Person
}

func (m *MemDb) DB() *memdb.MemDB {
	return m.db
}

func (m *MemDb) Insert(person people.Person) error {
	txn := m.db.Txn(true)
	defer txn.Abort()

	key := person.Nickname + " " + person.Name + " " + person.StackString()
	key = strings.ToLower(key)
	personMem := PersonMem{
		Key:    key,
		Person: person,
	}

	err := txn.Insert("people", personMem)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (m *MemDb) BulkInsert(people []people.Person) error {
	txn := m.db.Txn(true)
	defer txn.Abort()

	for _, person := range people {
		key := person.Nickname + " " + person.Name + " " + person.StackString()
		key = strings.ToLower(key)
		personMem := PersonMem{
			Key:    key,
			Person: person,
		}

		err := txn.Insert("people", personMem)
		if err != nil {
			return err
		}
	}

	txn.Commit()

	return nil
}

func (m *MemDb) Search(term string) ([]people.Person, error) {
	txn := m.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("people", "id")
	if err != nil {
		return nil, err
	}

	var people []people.Person

	for {
		if len(people) >= 50 {
			break
		}

		raw := it.Next()
		if raw == nil {
			break
		}

		personMem := raw.(PersonMem)
		if strings.Contains(personMem.Key, term) {
			people = append(people, personMem.Person)
		}
	}

	return people, nil
}

func (m *MemDb) Get(key string) (*people.Person, error) {
	txn := m.db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("people", "id", key)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, ErrMemDbNotFound
	}

	personMem := raw.(PersonMem)
	return &personMem.Person, nil
}

func NewMemDb() *MemDb {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"people": {
				Name: "people",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Key"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Fatalf("Failed to create memdb: %v", err)
	}

	return &MemDb{db: db}
}
