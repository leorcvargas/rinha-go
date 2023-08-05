package peopledb

import (
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"gorm.io/gorm"
)

type PersonRepository struct {
	db *gorm.DB
}

func (r *PersonRepository) Create(person *people.Person) (*people.Person, error) {
	model := NewPersonModel(
		person.Nickname,
		person.Name,
		person.Birthdate,
		person.Stack,
	)

	err := r.db.Create(&model).Error
	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func NewPersonRepository(db *gorm.DB) people.Repository {
	return &PersonRepository{db: db}
}
