package peopledb

import (
	"arena"
	"strings"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

type PersonItem struct {
	Key string
	people.Person
}

type Mem2 struct {
	list []PersonItem
}

func (m *Mem2) Add(person people.Person) {
	key := person.Nickname + " " + person.Name + " " + person.StackString()
	key = strings.ToLower(key)
	item := PersonItem{
		Key:    key,
		Person: person,
	}

	m.list = append(m.list, item)
}

func (m *Mem2) AddBatch(batch []people.Person) {
	a := arena.NewArena()
	defer a.Free()

	batchSize := len(batch)

	input := arena.MakeSlice[PersonItem](a, batchSize, batchSize)

	for i := 0; i < batchSize; i++ {
		item := batch[i]

		input[i] = PersonItem{
			Key:    item.Nickname + " " + item.Name + " " + item.StackString(),
			Person: item,
		}
	}

	m.list = append(m.list, input...)
}

func (m *Mem2) Search(query string) []people.Person {
	query = strings.ToLower(query)

	result := make([]people.Person, 0)

	size := len(m.list)
	limit := 50

	front := 0
	back := size - 1

	for i := 0; i < size; i++ {
		if len(result) >= limit {
			break
		}

		if strings.Contains(m.list[front].Key, query) {
			result = append(result, m.list[front].Person)
		}

		if strings.Contains(m.list[back].Key, query) {
			result = append(result, m.list[back].Person)
		}

		front++
		back--

		if front > back {
			break
		}
	}

	return result
}

func NewMem2() *Mem2 {
	return &Mem2{
		list: make([]PersonItem, 0),
	}
}
