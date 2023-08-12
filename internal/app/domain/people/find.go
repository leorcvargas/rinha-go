package people

type FindPeople struct {
	repository Repository
}

func (f *FindPeople) ByID(id string) (*Person, error) {
	person, err := f.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (f *FindPeople) Search(term string) ([]*Person, error) {
	people, err := f.repository.Search(term)
	if err != nil {
		return nil, err
	}

	return people, nil
}

func NewFindPeople(repository Repository) *FindPeople {
	return &FindPeople{repository: repository}
}
