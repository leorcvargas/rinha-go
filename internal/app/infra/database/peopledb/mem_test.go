package peopledb_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

func BenchmarkSearch_50000(b *testing.B) {
	sampleSize := 50000

	mem := peopledb.NewMem()
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

		mem.Add(*people.BuildPerson(
			fakePerson.ID,
			fakePerson.Nickname,
			fakePerson.Name,
			fakePerson.Birthdate,
			fakePerson.Stack,
		))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := mem.Search("go")
		if len(result) == 0 {
			b.Errorf("expected at least one result, got %v", result)
		}
	}
}

func BenchmarkSearch_100000(b *testing.B) {
	sampleSize := 100000

	mem := peopledb.NewMem()
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

		mem.Add(*people.BuildPerson(
			fakePerson.ID,
			fakePerson.Nickname,
			fakePerson.Name,
			fakePerson.Birthdate,
			fakePerson.Stack,
		))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := mem.Search("go")
		if len(result) == 0 {
			b.Errorf("expected at least one result, got %v", result)
		}
	}
}
