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
	peopleCache   *redis.Client
	nicknameCache *redis.Client
}

func (p *PeopleDbCache) Cache() *redis.Client {
	return p.peopleCache
}

func (p *PeopleDbCache) Get(key string) (*people.Person, error) {
	t := top("cache-get-nickname")
	defer t()

	item, err := p.peopleCache.Get(ctx, key).Result()
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
	t := top("cache-get-nickname")
	defer t()

	_, err := p.nicknameCache.Get(ctx, nickname).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *PeopleDbCache) Set(key string, person *people.Person) (*people.Person, error) {
	t := top("cache-set")
	defer t()

	errChan := make(chan error, 2)
	defer close(errChan)

	go func() {
		item, err := sonic.Marshal(person)
		if err != nil {
			errChan <- err
			return
		}

		_, err = p.peopleCache.Set(ctx, key, item, time.Hour).Result()
		if err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	go func() {
		err := p.SetNickname(person.Nickname)
		if err != nil {
			errChan <- err
			return
		}

		errChan <- nil
	}()

	for i := 0; i < 2; i++ {
		err := <-errChan
		if err != nil {
			return nil, err
		}
	}

	return person, nil
}

func (p *PeopleDbCache) SetNickname(nickname string) error {
	t := top("cache-set-nickname")
	defer t()

	_, err := p.nicknameCache.Set(ctx, nickname, true, time.Hour).Result()

	return err
}

func NewPeopleDbCache() *PeopleDbCache {
	peopleCache := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})
	nicknameCache := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       1,
	})

	return &PeopleDbCache{
		peopleCache:   peopleCache,
		nicknameCache: nicknameCache,
	}
}
