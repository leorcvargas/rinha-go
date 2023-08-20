package people

import (
	"strings"

	"github.com/google/uuid"
)

type Person struct {
	ID        string
	Nickname  string
	Name      string
	Birthdate string
	Stack     []string
}

func (p *Person) StackStr() string {
	return strings.Join(p.Stack, ",")
}

func (p *Person) SearchStr() string {
	return p.Nickname + " " + p.Name + " " + p.StackStr()
}

func BuildPerson(
	id string,
	nickname string,
	name string,
	birthdate string,
	stack []string,
) *Person {
	return &Person{
		ID:        id,
		Nickname:  nickname,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}

func NewPerson(
	nickname string,
	name string,
	birthdate string,
	stack []string,
) *Person {
	return &Person{
		ID:        uuid.NewString(),
		Nickname:  nickname,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}
