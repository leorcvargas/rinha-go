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
	t := top("cache-get-nickname")
	defer t()

	personBytes, err := p.client.Do(ctx, p.client.B().Get().Key(key).Build()).AsBytes()
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

	_, err := p.client.Do(ctx, p.client.B().Get().Key(nickname).Build()).AsBytes()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *PeopleDbCache) Set(key string, person *people.Person) (*people.Person, error) {
	t := top("cache-set")
	defer t()

	item, err := sonic.MarshalString(person)
	if err != nil {
		return nil, err
	}

	cmds := make(rueidis.Commands, 0, 2)
	cmds = append(cmds, p.client.B().Set().Key(person.ID).Value(item).Ex(time.Hour).Build())
	cmds = append(cmds, p.client.B().Set().Key(person.Nickname).Value("1").Ex(time.Hour).Build())

	for _, res := range p.client.DoMulti(ctx, cmds...) {
		err := res.Error()

		if err != nil {
			return nil, err
		}
	}

	return person, nil
}

func NewPeopleDbCache() *PeopleDbCache {
	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"cache:6379"}})
	if err != nil {
		panic(err)
	}

	return &PeopleDbCache{
		client: client,
	}
}
