package peopledb

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/redis/go-redis/v9"
)

type PersonRepository struct {
	db         *sql.DB
	cache      *PeopleDbCache
	insertChan InsertChan
}

func countOpTime(label string) func() {
	start := time.Now()

	return func() {
		elapsed := time.Since(start)
		log.Printf(">>> %s end %s", label, elapsed)
	}
}

func (p *PersonRepository) Create(person *people.Person) (*people.Person, error) {
	// strStack := strings.Join(person.Stack, ",")
	// searchAggregator := person.Nickname + person.Name + strStack
	// _, err := p.db.Exec(
	// 	"INSERT INTO people (id, nickname, name, birthdate, stack, search_aggr) VALUES ($1, $2, $3, $4, $5, $6)",
	// 	person.ID,
	// 	person.Nickname,
	// 	person.Name,
	// 	person.Birthdate,
	// 	strStack,
	// 	searchAggregator,
	// )
	// if err != nil {
	// 	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
	// 		return nil, people.ErrNicknameTaken
	// 	}

	// 	return nil, err
	// }

	// var count int
	// err := p.db.
	// 	QueryRow(
	// 		"SELECT COUNT(1) FROM people WHERE nickname = $1",
	// 		person.Nickname,
	// 	).
	// 	Scan(&count)
	// if err != nil {
	// 	return nil, err
	// }

	personSameNickname, err := p.cache.Get(person.Nickname)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if personSameNickname != nil {
		return nil, people.ErrNicknameTaken
	}

	// TODO refatorar pra salvar algo diferente em vez de repetir a mesma coisa
	x := countOpTime("create-person-repo")
	defer x()

	p.cache.Set(person.Nickname, person)

	p.cache.Set(person.ID.String(), person)

	p.insertChan <- *person

	return person, nil
}

func (p *PersonRepository) FindByID(id string) (*people.Person, error) {
	x1 := countOpTime("redis-get-person")
	cachedPerson, err := p.cache.Get(id)
	x1()

	if err != nil && err != redis.Nil {
		return nil, err
	}

	if cachedPerson != nil {
		return cachedPerson, nil
	}

	x2 := countOpTime("find-by-id-person")
	defer x2()
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
	x := countOpTime("search-person")
	defer x()
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

func NewPersonRepository(db *sql.DB, cache *PeopleDbCache, insertChan InsertChan) people.Repository {
	return &PersonRepository{
		db:         db,
		cache:      cache,
		insertChan: insertChan,
	}
}

type InsertChan chan people.Person

type InsertListenerRunner func(db *sql.DB) InsertChan

func runInsertListener(db *sql.DB) InsertChan {
	var insertChan = make(InsertChan)

	go insertWorker(insertChan, db)

	return insertChan
}

// isso aqui é um crime de guerra
func insertWorker(ch InsertChan, db *sql.DB) {
	// mem := arena.NewArena()

	// items := arena.MakeSlice[*people.Person](mem, 50, 2048)
	// at := 0
	// lastInsertedAt := time.Time{}

	items := make([]*people.Person, 0, 50)

	for {
		select {
		case item := <-ch:
			items = append(items, &item)

			if len(items) >= 50 {
				saveBatch(items, db)
				items = make([]*people.Person, 0, 50)
			}
		case <-time.After(5 * time.Second):
			if len(items) > 0 {
				saveBatch(items, db)
				items = make([]*people.Person, 0, 50)
			}
		}
	}
}

func saveBatch(items []*people.Person, db *sql.DB) {
	valueStrings := make([]string, 0, len(items))
	valueArgs := make([]interface{}, 0, len(items)*6)
	for i, item := range items {
		strStack := strings.Join(item.Stack, ",")

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		valueArgs = append(valueArgs, item.ID)
		valueArgs = append(valueArgs, item.Nickname)
		valueArgs = append(valueArgs, item.Name)
		valueArgs = append(valueArgs, item.Birthdate)
		valueArgs = append(valueArgs, strStack)
	}

	statement := fmt.Sprintf(
		"INSERT INTO people (id, nickname, name, birthdate, stack) VALUES %s",
		strings.Join(valueStrings, ","),
	)
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Printf("failed to prepare statement %v", err)
		return
	}
	defer stmt.Close()

	stmt.Exec(valueArgs...)
}
