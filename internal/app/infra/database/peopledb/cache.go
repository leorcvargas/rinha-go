package peopledb

import (
	"context"
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
	t := top("cache-get")
	defer t()

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
	t := top("cache-get-nickname")
	defer t()

	getNicknameCmd := p.client.
		B().
		Getbit().
		Key("nickname:" + nickname).
		Offset(0).
		Cache()

	return p.client.DoCache(ctx, getNicknameCmd, time.Hour).AsBool()
}

func (p *PeopleDbCache) Set(key string, person *people.Person) (*people.Person, error) {
	t := top("cache-set")
	defer t()

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

func NewPeopleDbCache() *PeopleDbCache {
	opts := rueidis.ClientOption{
		InitAddress: []string{"cache:6379"},
	}
	client, err := rueidis.NewClient(opts)
	if err != nil {
		panic(err)
	}

	return &PeopleDbCache{
		client: client,
	}
}
