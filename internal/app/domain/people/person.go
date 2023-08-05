package people

type Person struct {
	ID        string   `json:"id"`
	Nickname  string   `json:"apelido"`
	Name      string   `json:"nome"`
	Birthdate string   `json:"nascimento"`
	Stack     []string `json:"stack"`
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
		Nickname:  nickname,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}
