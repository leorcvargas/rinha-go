package peopledb

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type PeopleDbCache struct {
	cache *redis.Client
}

func (p *PeopleDbCache) Cache() *redis.Client {
	return p.cache
}

func (p *PeopleDbCache) Get(key string) (*people.Person, error) {
	// x := time.Now()
	// defer func() {
	// 	log.Info("Get cache time: ", time.Since(x))
	// }()

	item, err := p.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var person people.Person
	err = sonic.Unmarshal([]byte(item), &person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (p *PeopleDbCache) GetNickname(nickname string) (bool, error) {
	// x := time.Now()
	// defer func() {
	// 	log.Info("GetNickname cache time: ", time.Since(x))
	// }()

	_, err := p.cache.Get(ctx, nickname).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *PeopleDbCache) Set(key string, person *people.Person) (*people.Person, error) {
	// x := time.Now()
	// defer func() {
	// 	log.Info("Set cache time: ", time.Since(x))
	// }()

	item, err := sonic.Marshal(person)
	if err != nil {
		return nil, err
	}

	_, err = p.cache.Set(ctx, key, item, time.Hour).Result()
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p *PeopleDbCache) SetNickname(nickname string) error {
	// x := time.Now()
	// defer func() {
	// 	log.Info("SetNickname cache time: ", time.Since(x))
	// }()

	_, err := p.cache.Set(ctx, nickname, true, time.Hour).Result()

	return err
}

func NewPeopleDbCache() *PeopleDbCache {
	cache := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})

	return &PeopleDbCache{
		cache: cache,
	}
}
