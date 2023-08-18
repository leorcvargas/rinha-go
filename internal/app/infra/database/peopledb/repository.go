package peopledb

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/redis/rueidis"
)

type PersonRepository struct {
	db         *pgxpool.Pool
	cache      *PeopleDbCache
	insertChan chan people.Person
	mem2       *Mem2
}

func (p *PersonRepository) Create(person *people.Person) (*people.Person, error) {
	_, err := p.db.Exec(
		context.Background(),
		InsertPersonQuery,
		person.ID,
		person.Nickname,
		person.Name,
		person.Birthdate,
		person.StackString(),
		strings.ToLower(person.Nickname+" "+person.Name+" "+person.StackString()),
	)

	if err != nil {
		return nil, err
	}

	go p.cache.Set(person.ID, person)

	return person, nil

	// nicknameTaken, err := p.cache.GetNickname(person.Nickname)
	// if err != nil {
	// 	return nil, err
	// }

	// if nicknameTaken {
	// 	return nil, people.ErrNicknameTaken
	// }

	// p.cache.Set(person.ID, person)

	// p.insertChan <- *person

	// return person, nil
}

func (p *PersonRepository) FindByID(id string) (*people.Person, error) {
	cachedPerson, err := p.cache.Get(id)

	if err != nil && !rueidis.IsRedisNil(err) {
		return nil, err
	}

	if cachedPerson != nil {
		return cachedPerson, nil
	}

	var person people.Person
	var strStack string

	err = p.db.QueryRow(
		context.Background(),
		SelectPersonByIDQuery,
		id,
	).Scan(
		&person.ID,
		&person.Nickname,
		&person.Name,
		&person.Birthdate,
		&strStack,
	)
	person.Stack = strings.Split(strStack, ",")

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, people.ErrPersonNotFound
		}

		return nil, err
	}

	return &person, nil
}

// func (p *PersonRepository) Search(term string) ([]people.Person, error) {
// 	result := p.mem2.Search(term)
// 	return result, nil
// }

func (p *PersonRepository) Search(term string) ([]people.Person, error) {
	return p.searchTrigram(term)
}

func (p *PersonRepository) searchTrigram(term string) ([]people.Person, error) {
	rows, err := p.db.Query(
		context.Background(),
		SearchPeopleTrgmQuery,
		term,
	)
	if err != nil {
		log.Errorf("Error executing trigram search: %v", err)
		return nil, err
	}

	return mapSearchResult(rows)
}

func (p *PersonRepository) CountAll() (int64, error) {
	var total int64

	err := p.db.
		QueryRow(
			context.Background(),
			CountPeopleQuery,
		).
		Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func mapSearchResult(rows pgx.Rows) ([]people.Person, error) {
	result := make([]people.Person, 0)
	for rows.Next() {
		var person people.Person
		var strStack string
		var birthdate time.Time

		err := rows.Scan(
			&person.ID,
			&person.Nickname,
			&person.Name,
			&birthdate,
			&strStack,
		)
		if err != nil {
			log.Errorf("Error scanning row: %v", err)
			return nil, err
		}

		person.Stack = strings.Split(strStack, ",")
		person.Birthdate = birthdate.Format("2006-01-02")

		result = append(result, person)
	}

	return result, nil
}

func NewPersonRepository(db *pgxpool.Pool, cache *PeopleDbCache, mem2 *Mem2, insertChan chan people.Person) people.Repository {
	return &PersonRepository{
		db:         db,
		cache:      cache,
		insertChan: insertChan,
		mem2:       mem2,
	}
}
