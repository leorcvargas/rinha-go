package people

type CreatePeople struct {
	repository Repository
}

func (c *CreatePeople) Execute(person *Person) (*Person, error) {
	return c.repository.Create(person)
}

func NewCreatePeople(repository Repository) *CreatePeople {
	return &CreatePeople{repository: repository}
}
