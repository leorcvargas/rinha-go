package people

type Repository interface {
	Create(person *Person) (*Person, error)
}
