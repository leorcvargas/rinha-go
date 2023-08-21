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

type PeopleDbCache struct {
	client rueidis.Client
}

func (p *PeopleDbCache) Cache() rueidis.Client {
	return p.client
}

func (p *PeopleDbCache) Get(key string) (*people.Person, error) {
	getCmd := p.client.
		B().
		Get().
		Key("person:" + key).
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

func (p *PeopleDbCache) GetNickname(nickname string) (bool, error) {
	getNicknameCmd := p.client.
		B().
		Getbit().
		Key("nickname:" + nickname).
		Offset(0).
		Cache()

	return p.client.DoCache(ctx, getNicknameCmd, time.Hour).AsBool()
}

func (p *PeopleDbCache) Set(key string, person *people.Person) (*people.Person, error) {
	item, err := sonic.MarshalString(person)
	if err != nil {
		return nil, err
	}

	setPersonCmd := p.client.
		B().
		Set().
		Key("person:" + person.ID).
		Value(item).
		Ex(time.Hour).
		Build()

	setNicknameCmd := p.client.
		B().
		Setbit().
		Key("nickname:" + person.Nickname).
		Offset(0).
		Value(1).
		Build()

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds, setPersonCmd)
	cmds = append(cmds, setNicknameCmd)

	for _, res := range p.client.DoMulti(ctx, cmds...) {
		err := res.Error()

		if err != nil {
			return nil, err
		}
	}

	return person, nil
}

func (p *PeopleDbCache) SetSearch(term string, result []people.Person) error {
	item, err := sonic.MarshalString(result)
	if err != nil {
		return err
	}

	setSearchCmd := p.client.
		B().
		Set().
		Key("search:" + term).
		Value(item).
		Ex(30 * time.Second).
		Build()

	return p.client.Do(ctx, setSearchCmd).Error()
}

func (p *PeopleDbCache) GetSearch(term string) ([]people.Person, error) {
	getSearchCmd := p.client.
		B().
		Get().
		Key("search:" + term).
		Cache()

	resultBytes, err := p.client.
		DoCache(
			ctx,
			getSearchCmd,
			30*time.Second,
		).
		AsBytes()

	if err != nil {
		return nil, err
	}

	var result []people.Person
	err = sonic.Unmarshal(resultBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewPeopleDbCache() *PeopleDbCache {
	address := fmt.Sprintf(
		"%s:%s",
		os.Getenv("CACHE_HOST"),
		os.Getenv("CACHE_PORT"),
	)

	opts := rueidis.ClientOption{
		InitAddress: []string{address},
	}
	client, err := rueidis.NewClient(opts)
	if err != nil {
		panic(err)
	}

	return &PeopleDbCache{
		client: client,
	}
}
