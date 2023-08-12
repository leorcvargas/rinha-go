package peopledb

import (
	"database/sql"
	"strings"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type PersonRepository struct {
	db    *sql.DB
	cache *PeopleDbCache
}

func (p *PersonRepository) Create(person *people.Person) (*people.Person, error) {
	strStack := strings.Join(person.Stack, ",")
	_, err := p.db.Exec(
		"INSERT INTO people (id, nickname, name, birthdate, stack) VALUES ($1, $2, $3, $4, $5)",
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
		"SELECT id, nickname, name, birthdate, stack FROM people WHERE id = $1",
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

func (p *PersonRepository) Search(term string) ([]*people.Person, error) {
	rows, err := p.db.Query(
		`SELECT id, nickname, name, birthdate, stack FROM people p
		WHERE p.fts_q @@ websearch_to_tsquery('english', $1)
		LIMIT 50;`,
		term,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]*people.Person, 0)
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

		result = append(result, &person)
	}

	return result, nil
}

func (p *PersonRepository) CountAll() (int64, error) {
	var total int64

	err := p.db.QueryRow(
		"SELECT COUNT(*) FROM people",
	).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func NewPersonRepository(db *sql.DB, cache *PeopleDbCache) people.Repository {
	return &PersonRepository{
		db:    db,
		cache: cache,
	}
}
