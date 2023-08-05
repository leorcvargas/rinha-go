package peopledb

import (
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type PersonModel struct {
	gorm.Model
	ID        string `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Nickname  string `gorm:"type:varchar(32)"`
	Name      string `gorm:"type:varchar(100)"`
	Birthdate string
	Stack     pq.StringArray `gorm:"type:text[]"`
}

func (PersonModel) TableName() string {
	return "people"
}

func (p *PersonModel) ToDomain() *people.Person {
	return people.BuildPerson(
		p.ID,
		p.Nickname,
		p.Name,
		p.Birthdate,
		p.Stack,
	)
}

func BuildPersonModel(id, nickname, name, birthdate string, stack []string) *PersonModel {
	return &PersonModel{
		ID:        id,
		Nickname:  nickname,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}

func NewPersonModel(nickname, name, birthdate string, stack []string) *PersonModel {
	return &PersonModel{
		Nickname:  nickname,
		Name:      name,
		Birthdate: birthdate,
		Stack:     stack,
	}
}
