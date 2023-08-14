package peopledb

import (
	"strings"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

type PeopleMemoryStorage struct {
	list []people.Person
	size uint64
}

func (p *PeopleMemoryStorage) Insert(person people.Person) (uint64, error) {
	p.list = append(p.list, person)
	p.size++

	return p.size, nil
}

func (p *PeopleMemoryStorage) Search(term string) []people.Person {
	limit := 50
	result := make([]people.Person, 0, limit)

	for _, person := range p.list {
		search := person.Name + " " + person.Nickname + " " + strings.Join(person.Stack, " ")

		if strings.Contains(search, term) {
			result = append(result, person)
		}

		if len(result) >= limit {
			break
		}
	}

	return result
}

func NewPeopleMemoryStorage() *PeopleMemoryStorage {
	return &PeopleMemoryStorage{
		list: make([]people.Person, 0, 100000),
		size: 0,
	}
}
