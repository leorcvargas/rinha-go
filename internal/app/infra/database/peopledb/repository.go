package peopledb

import (
	"database/sql"
	"strings"
	"sync"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/redis/go-redis/v9"
)

type PersonRepository struct {
	db         *sql.DB
	cache      *PeopleDbCache
	insertChan chan people.Person
	mem2       *Mem2
}

func (p *PersonRepository) Create(person *people.Person) (*people.Person, error) {
	nicknameTaken, err := p.cache.GetNickname(person.Nickname)
	if err != nil {
		return nil, err
	}

	if nicknameTaken {
		return nil, people.ErrNicknameTaken
	}

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		p.cache.Set(person.ID, person)
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		p.cache.SetNickname(person.Nickname)
		wg.Done()
	}()

	p.insertChan <- *person

	wg.Wait()

	return person, nil
}

func (p *PersonRepository) FindByID(id string) (*people.Person, error) {
	cachedPerson, err := p.cache.Get(id)

	if err != nil && err != redis.Nil {
		return nil, err
	}

	if cachedPerson != nil {
		return cachedPerson, nil
	}

	var person people.Person
	var strStack string

	err = p.db.QueryRow(
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
		if err == sql.ErrNoRows {
			return nil, people.ErrPersonNotFound
		}

		return nil, err
	}

	return &person, nil
}

func (p *PersonRepository) Search(term string) ([]people.Person, error) {
	result := p.mem2.Search(term)
	return result, nil
}

func (p *PersonRepository) CountAll() (int64, error) {
	var total int64

	err := p.db.QueryRow(
		CountPeopleQuery,
	).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func mapSearchResult(rows *sql.Rows) ([]people.Person, error) {
	result := make([]people.Person, 0)
	for rows.Next() {
		var person people.Person
		var strStack string
		var birthdate string

		err := rows.Scan(
			&person.ID,
			&person.Nickname,
			&person.Name,
			&birthdate,
			&strStack,
		)
		if err != nil {
			return nil, err
		}

		person.Stack = strings.Split(strStack, ",")
		person.Birthdate = birthdate[0:10]

		result = append(result, person)
	}

	return result, nil
}

func NewPersonRepository(db *sql.DB, cache *PeopleDbCache, mem2 *Mem2, insertChan chan people.Person) people.Repository {
	return &PersonRepository{
		db:         db,
		cache:      cache,
		insertChan: insertChan,
		mem2:       mem2,
	}
}
