package pubsub

import (
	"context"
	"encoding/json"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

type PersonInsertSubscriber struct {
	cache *peopledb.PeopleDbCache
	memDb *peopledb.MemDb
}

func (p *PersonInsertSubscriber) Subscribe() {
	sub := p.cache.Cache().Subscribe(context.Background(), EventPersonInsert)
	defer sub.Close()

	ch := sub.Channel()

	for msg := range ch {
		var people []people.Person

		err := json.Unmarshal([]byte(msg.Payload), &people)
		if err != nil {
			panic(err)
		}

		for _, person := range people {
			p.memDb.Insert(person)
		}
	}
}

func NewPersonInsertSubscriber(cache *peopledb.PeopleDbCache, memDb *peopledb.MemDb) *PersonInsertSubscriber {
	return &PersonInsertSubscriber{
		cache: cache,
		memDb: memDb,
	}
}
