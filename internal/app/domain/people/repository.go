package people

type Repository interface {
	Create(person *Person) error
	CheckNicknameExists(nickname string) (bool, error)
	FindByID(id string) (*Person, error)
	Search(term string) ([]Person, error)
	CountAll() (int64, error)
}
