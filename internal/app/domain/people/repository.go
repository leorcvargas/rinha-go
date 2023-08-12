package people

type Repository interface {
	Create(person *Person) (*Person, error)
	FindByID(id string) (*Person, error)
	Search(term string) ([]*Person, error)
	CountAll() (int64, error)
}
