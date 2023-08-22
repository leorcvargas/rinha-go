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
	client rueidis.Client
}

func (p *Cache) Get(key string) (*people.Person, error) {
	getCmd := p.client.
		B().
		Get().
		Key(key).
		Cache()

	personBytes, err := p.client.DoCache(ctx, getCmd, time.Hour).AsBytes()
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
	getNicknameCmd := p.client.
		B().
		Getbit().
		Key(nickname).
		Offset(0).
		Cache()

	return p.client.DoCache(ctx, getNicknameCmd, time.Hour).AsBool()
}

func (p *Cache) Set(person *people.Person) error {
	item, err := sonic.MarshalString(person)
	if err != nil {
		return err
	}

	setPersonCmd := p.client.
		B().
		Set().
		Key(person.ID).
		Value(item).
		Ex(15 * time.Second).
		Build()

	setNicknameCmd := p.client.
		B().
		Setbit().
		Key(person.Nickname).
		Offset(0).
		Value(1).
		Build()

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds, setPersonCmd)
	cmds = append(cmds, setNicknameCmd)

	for _, res := range p.client.DoMulti(ctx, cmds...) {
		err := res.Error()

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
	client, err := rueidis.NewClient(opts)
	if err != nil {
		panic(err)
	}

	return &Cache{
		client: client,
	}
}
