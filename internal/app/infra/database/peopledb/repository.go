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
	db       *pgxpool.Pool
	cache    *PeopleDbCache
	jobQueue JobQueue
}

func (p *PersonRepository) Create(person *people.Person) error {
	if _, err := p.cache.Set(person.ID, person); err != nil {
		log.Errorf("Error setting person in cache: %v", err)
		return err
	}

	p.jobQueue <- Job{Payload: person}

	return nil
}

func (p *PersonRepository) FindByID(id string) (*people.Person, error) {
	cachedPerson, err := p.cache.Get(id)

	if err != nil && !rueidis.IsRedisNil(err) {
		log.Errorf("Error getting person from cache: %v", err)
		return nil, err
	}

	if cachedPerson != nil {
		return cachedPerson, nil
	}

	var person people.Person
	var strStack string
	var birthdate time.Time

	err = p.db.QueryRow(
		context.Background(),
		SelectPersonByIDQuery,
		id,
	).Scan(
		&person.ID,
		&person.Nickname,
		&person.Name,
		&birthdate,
		&strStack,
	)
	person.Stack = strings.Split(strStack, ",")
	person.Birthdate = birthdate.Format("2006-01-02")

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, people.ErrPersonNotFound
		}

		log.Errorf("Error querying person: %v", err)

		return nil, err
	}

	return &person, nil
}

func (p *PersonRepository) Search(term string) ([]people.Person, error) {
	sanitizedTerm := strings.ToLower(term)

	cachedResult, err := p.cache.GetSearch(sanitizedTerm)
	if err != nil && !rueidis.IsRedisNil(err) {
		log.Errorf("Error getting search from cache: %v", err)
		return nil, err
	}

	if len(cachedResult) > 0 {
		log.Infof("Returning cached search result for term: %s", sanitizedTerm)
		return cachedResult, nil
	}

	rows, err := p.db.Query(
		context.Background(),
		SearchPeopleTrgmQuery,
		sanitizedTerm,
	)
	if err != nil {
		log.Errorf("Error executing trigram search: %v", err)
		return nil, err
	}

	result, err := p.mapSearchResult(rows)
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		go func() {
			if err := p.cache.SetSearch(sanitizedTerm, result); err != nil {
				log.Errorf("Error setting search result in cache: %v", err)
			}
		}()
	}

	return result, nil
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
		log.Errorf("Error counting people: %v", err)
		return 0, err
	}

	return total, nil
}

func (p *PersonRepository) CheckNicknameExists(nickname string) (bool, error) {
	nicknameTaken, err := p.cache.GetNickname(nickname)
	if err != nil {
		return false, err
	}

	return nicknameTaken, nil
}

func (p *PersonRepository) mapSearchResult(rows pgx.Rows) ([]people.Person, error) {
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

func NewPersonRepository(db *pgxpool.Pool, cache *PeopleDbCache, jobQueue JobQueue) people.Repository {
	return &PersonRepository{
		db:       db,
		cache:    cache,
		jobQueue: jobQueue,
	}
}
