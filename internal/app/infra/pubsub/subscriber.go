package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
		now := time.Now()
		// memory := arena.NewArena()

		// people := arena.MakeSlice[people.Person](memory, 0, 50)

		var people []people.Person

		err := json.Unmarshal([]byte(msg.Payload), &people)
		if err != nil {
			panic(err)
		}

		p.memDb.BulkInsert(people)
		log.Printf("memdb synced in %s, message size %d", time.Since(now), len(people))
	}
}

func NewPersonInsertSubscriber(cache *peopledb.PeopleDbCache, memDb *peopledb.MemDb) *PersonInsertSubscriber {
	return &PersonInsertSubscriber{
		cache: cache,
		memDb: memDb,
	}
}
