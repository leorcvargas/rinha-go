package worker

import "github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"

func CreateInsertChannel() chan people.Person {
	return make(chan people.Person)
}
