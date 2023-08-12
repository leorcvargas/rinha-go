package people

import "github.com/google/uuid"

type Person struct {
	ID        uuid.UUID
	Nickname  string
	Name      string
	Birthdate string
	Stack     []string
}

func BuildPerson(
	id uuid.UUID,
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
		ID:        uuid.New(),
		Nickname:  nickname,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}
