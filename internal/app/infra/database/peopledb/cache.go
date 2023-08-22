package peopledb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/redis/rueidis"
)

var ctx = context.Background()

type Cache struct {
	peopleClient   rueidis.Client
	nicknameClient rueidis.Client
}

func (p *Cache) Get(key string) (*people.Person, error) {
	getCmd := p.peopleClient.
		B().
		Get().
		Key(key).
		Cache()

	personBytes, err := p.peopleClient.DoCache(ctx, getCmd, time.Hour).AsBytes()
	if err != nil {
		return nil, err
	}

	var person people.Person
	err = sonic.Unmarshal(personBytes, &person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

func (p *Cache) GetNickname(nickname string) (bool, error) {
	getNicknameCmd := p.nicknameClient.
		B().
		Getbit().
		Key(nickname).
		Offset(0).
		Cache()

	return p.nicknameClient.DoCache(ctx, getNicknameCmd, time.Hour).AsBool()
}

func (p *Cache) Set(person *people.Person) error {
	errorChannel := make(chan error, 2)

	go func() {
		item, marshalError := sonic.MarshalString(person)
		if marshalError != nil {
			errorChannel <- marshalError
			return
		}

		cmd := p.peopleClient.
			B().
			Set().
			Key(person.ID).
			Value(item).
			Ex(15 * time.Minute).
			Build()

		setError := p.peopleClient.Do(ctx, cmd).Error()
		if setError != nil {
			errorChannel <- setError
			return
		}

		errorChannel <- nil
	}()

	go func() {
		cmd := p.nicknameClient.
			B().
			Setbit().
			Key(person.Nickname).
			Offset(0).
			Value(1).
			Build()

		setError := p.nicknameClient.Do(ctx, cmd).Error()
		if setError != nil {
			errorChannel <- setError
			return
		}

		errorChannel <- nil
	}()

	for i := 0; i < 2; i++ {
		err := <-errorChannel

		if err != nil {
			return err
		}
	}

	return nil
}

func NewCache() *Cache {
	address := fmt.Sprintf(
		"%s:%s",
		os.Getenv("CACHE_HOST"),
		os.Getenv("CACHE_PORT"),
	)

	opts := rueidis.ClientOption{
		InitAddress:      []string{address},
		AlwaysPipelining: true,
		SelectDB:         0,
	}
	peopleClient, err := rueidis.NewClient(opts)
	if err != nil {
		panic(err)
	}

	opts = rueidis.ClientOption{
		InitAddress:      []string{address},
		AlwaysPipelining: true,
		SelectDB:         1,
	}
	nicknameClient, err := rueidis.NewClient(opts)
	if err != nil {
		panic(err)
	}

	return &Cache{
		peopleClient:   peopleClient,
		nicknameClient: nicknameClient,
	}
}
