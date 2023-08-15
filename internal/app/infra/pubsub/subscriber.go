package pubsub

import (
	"context"
	"encoding/json"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

type PersonInsertSubscriber struct {
	cache *peopledb.PeopleDbCache
	mem2  *peopledb.Mem2
}

func (p *PersonInsertSubscriber) Subscribe() {
	sub := p.cache.Cache().Subscribe(context.Background(), EventPersonInsert)
	defer sub.Close()

	ch := sub.Channel()

	for msg := range ch {
		// memory := arena.NewArena()

		// people := arena.MakeSlice[people.Person](memory, 0, 50)

		var people []people.Person

		err := json.Unmarshal([]byte(msg.Payload), &people)
		if err != nil {
			panic(err)
		}

		p.mem2.AddBatch(people)
	}
}

func NewPersonInsertSubscriber(cache *peopledb.PeopleDbCache, mem2 *peopledb.Mem2) *PersonInsertSubscriber {
	return &PersonInsertSubscriber{
		cache: cache,
		mem2:  mem2,
	}
}
