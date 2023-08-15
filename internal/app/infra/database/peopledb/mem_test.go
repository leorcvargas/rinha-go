package peopledb_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

func BenchmarkMem2Search(b *testing.B) {
	sampleSize := 50000

	mem2 := peopledb.NewMem2()
	for i := 0; i < sampleSize; i++ {
		var stack []string

		if i%150 == 0 {
			stack = []string{"golang", "nodejs"}
		} else {
			stack = []string{faker.Word(), faker.Word()}
		}

		fakePerson := people.NewPerson(
			faker.Username(),
			faker.Name(),
			faker.Date(),
			stack,
		)
		err := faker.FakeData(&fakePerson)
		if err != nil {
			b.Errorf("expected no errors, got %v", err)
		}

		mem2.Add(*people.BuildPerson(
			fakePerson.ID,
			fakePerson.Nickname,
			fakePerson.Name,
			fakePerson.Birthdate,
			fakePerson.Stack,
		))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := mem2.Search("go")
		if len(result) == 0 {
			b.Errorf("expected at least one result, got %v", result)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	sampleSize := 50000

	var fakePeople []people.Person

	for i := 0; i < sampleSize; i++ {
		var stack []string

		if i%150 == 0 {
			stack = []string{"golang", "nodejs"}
		} else {
			stack = []string{faker.Word(), faker.Word()}
		}

		fakePerson := people.NewPerson(
			faker.Username(),
			faker.Name(),
			faker.Date(),
			stack,
		)
		err := faker.FakeData(&fakePerson)
		if err != nil {
			b.Errorf("expected no errors, got %v", err)
		}

		fakePeople = append(fakePeople, *fakePerson)
	}

	memDb := peopledb.NewMemDb()
	for _, person := range fakePeople {
		err := memDb.Insert(*people.BuildPerson(person.ID, person.Nickname, person.Name, person.Birthdate, person.Stack))
		if err != nil {
			b.Errorf("expected no errors, got %v", err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := memDb.Search("go")
		if err != nil {
			b.Errorf("expected no errors, got %v", err)
		}
	}
}
