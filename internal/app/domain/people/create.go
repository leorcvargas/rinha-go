package people

type CreatePerson struct {
	repository Repository
}

func (c *CreatePerson) Execute(
	Nickname string,
	Name string,
	Birthdate string,
	Stack []string,
) (*Person, error) {
	person := NewPerson(Nickname, Name, Birthdate, Stack)

	_, err := c.repository.Create(person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func NewCreatePerson(repository Repository) *CreatePerson {
	return &CreatePerson{repository: repository}
}
