package pubsub

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/log"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	"github.com/redis/rueidis"
)

type PersonInsertSubscriber struct {
	cache *peopledb.PeopleDbCache
	mem2  *peopledb.Mem2
}

func (p *PersonInsertSubscriber) handle(msg rueidis.PubSubMessage) {
	log.Infof("received %v", msg)
	var people []people.Person

	err := sonic.Unmarshal([]byte(msg.Message), &people)
	if err != nil {
		log.Errorf("Error on unmarshal message %s - err %v", msg.Message, err)
		panic(err)
	}

	p.mem2.AddBatch(people)
}

func (p *PersonInsertSubscriber) Subscribe() {
	err := p.cache.Cache().Receive(
		context.Background(),
		p.cache.Cache().B().Subscribe().Channel(EventPersonInsert).Build(),
		p.handle,
	)

	if err != nil {
		log.Errorf("Error on subscribe to %s - err %v", EventPersonInsert, err)
		panic(err)
	}
}

func NewPersonInsertSubscriber(cache *peopledb.PeopleDbCache, mem2 *peopledb.Mem2) *PersonInsertSubscriber {
	return &PersonInsertSubscriber{
		cache: cache,
		mem2:  mem2,
	}
}
