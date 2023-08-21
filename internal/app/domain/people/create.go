package people

type CreatePerson struct {
	repository Repository
}

func (c *CreatePerson) Execute(nickname string, name string, birthdate string, stack []string) (*Person, error) {
	person := NewPerson(nickname, name, birthdate, stack)

	err := c.repository.Create(person)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func NewCreatePerson(repository Repository) *CreatePerson {
	return &CreatePerson{repository: repository}
}
