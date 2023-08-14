package peopledb

import (
	"database/sql"
	"strings"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type PersonRepository struct {
	db               *sql.DB
	cache            *PeopleDbCache
	roundRobinSearch *roundRobinSearch
	searchdb         *sql.DB
}

func (p *PersonRepository) Create(person *people.Person) (*people.Person, error) {
	strStack := strings.Join(person.Stack, ",")
	_, err := p.db.Exec(
		InsertPersonQuery,
		person.ID,
		person.Nickname,
		person.Name,
		person.Birthdate,
		strStack,
	)
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
			return nil, people.ErrNicknameTaken
		}

		return nil, err
	}

	p.cache.Set(person.ID.String(), person)

	go func() {
		p.searchdb.Exec(
			InsertPersonQuery,
			person.ID,
			person.Nickname,
			person.Name,
			person.Birthdate,
			strStack,
		)
	}()

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
	result, err := p.searchFts(term)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}

	return p.roundRobinSearch.Search(term)
}

func (p *PersonRepository) searchFts(term string) ([]people.Person, error) {
	rows, err := p.db.Query(
		SearchPeopleFtsQuery,
		term,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result, err := mapSearchResult(rows)
	if err != nil {
		return nil, err
	}

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

		err := rows.Scan(
			&person.ID,
			&person.Nickname,
			&person.Name,
			&person.Birthdate,
			&strStack,
		)
		person.Stack = strings.Split(strStack, ",")
		if err != nil {
			return nil, err
		}

		result = append(result, person)
	}

	return result, nil
}

func NewPersonRepository(db *sql.DB, cache *PeopleDbCache) people.Repository {
	sqlitedb := NewSearchDatabase()

	searchers := []TrigramSearcher{
		&LocalTrigramSearcher{db: sqlitedb},
		&PsqlTrigramSearcher{db: sqlitedb},
	}

	rrs := NewRoundRobinSearcher(searchers...)

	return &PersonRepository{
		db:               db,
		cache:            cache,
		roundRobinSearch: rrs,
		searchdb:         sqlitedb,
	}
}
