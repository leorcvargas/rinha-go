package people

type CountPeople struct {
	repository Repository
}

func (c *CountPeople) CountAll() (int64, error) {
	total, err := c.repository.CountAll()
	if err != nil {
		return 0, err
	}

	return total, nil
}

func NewCountPeople(repository Repository) *CountPeople {
	return &CountPeople{repository: repository}
}
