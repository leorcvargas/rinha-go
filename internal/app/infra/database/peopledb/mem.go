package peopledb

import (
	"arena"
	"strings"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

const maxMemItems = 100000

type PersonItem struct {
	Key string
	people.Person
}

type Mem struct {
	list []PersonItem
}

func (m *Mem) Add(person people.Person) {
	key := person.Nickname + " " + person.Name + " " + person.StackString()
	key = strings.ToLower(key)
	item := PersonItem{
		Key:    key,
		Person: person,
	}

	m.list = append(m.list, item)
}

func (m *Mem) AddBatch(batch []people.Person) {
	a := arena.NewArena()
	defer a.Free()

	batchSize := len(batch)

	input := arena.MakeSlice[PersonItem](a, batchSize, batchSize)

	for i := 0; i < batchSize; i++ {
		item := batch[i]

		key := strings.ToLower(item.Nickname + item.Name + strings.Join(item.Stack, ""))

		input[i] = PersonItem{
			Key:    key,
			Person: item,
		}
	}

	m.list = append(m.list, arena.Clone(input)...)
}

func (m *Mem) Search(query string) []people.Person {
	query = strings.ToLower(query)

	limit := 50
	size := len(m.list)
	result := make([]people.Person, 0, limit)

	if size == 1 {
		if strings.Contains(m.list[0].Key, query) {
			result = append(result, m.list[0].Person)
		}
		return result
	}

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

func NewMem() *Mem {
	mem := &Mem{
		list: make([]PersonItem, 0, maxMemItems),
	}

	return mem
}
