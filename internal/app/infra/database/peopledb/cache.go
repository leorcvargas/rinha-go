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
		Key("person:" + key).
		Build()

	personBytes, err := p.client.Do(ctx, getCmd).AsBytes()
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
		Key("nickname:" + nickname).
		Offset(0).
		Build()

	return p.client.Do(ctx, getNicknameCmd).AsBool()
}

func (p *Cache) Set(person *people.Person) error {
	item, err := sonic.MarshalString(person)
	if err != nil {
		return err
	}

	setPersonCmd := p.client.
		B().
		Set().
		Key("person:" + person.ID).
		Value(item).
		Ex(time.Minute).
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
			return err
		}
	}

	return nil
}

func (p *Cache) SetSearch(term string, result []people.Person) error {
	item, err := sonic.MarshalString(result)
	if err != nil {
		return err
	}

	setSearchCmd := p.client.
		B().
		Set().
		Key("search:" + term).
		Value(item).
		Ex(1.5 * 60000 * time.Millisecond).
		Build()

	return p.client.Do(ctx, setSearchCmd).Error()
}

func (p *Cache) GetSearch(term string) ([]people.Person, error) {
	getSearchCmd := p.client.
		B().
		Get().
		Key("search:" + term).
		Build()

	resultBytes, err := p.client.
		Do(
			ctx,
			getSearchCmd,
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

func NewCache() *Cache {
	address := fmt.Sprintf(
		"%s:%s",
		os.Getenv("CACHE_HOST"),
		os.Getenv("CACHE_PORT"),
	)

	opts := rueidis.ClientOption{
		InitAddress: []string{address},
		// AlwaysPipelining: true,
		// CacheSizeEachConn: 256 * (1 << 20),
		// PipelineMultiplex: 8,
	}
	client, err := rueidis.NewClient(opts)
	if err != nil {
		panic(err)
	}

	return &Cache{client: client}
}
