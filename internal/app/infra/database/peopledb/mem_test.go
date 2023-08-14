package peopledb_test

import (
	"log"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
)

// func TestNewMemDb(t *testing.T) {
// 	memDb := peopledb.NewMemDb()
// 	if memDb == nil {
// 		t.Error("MemDb not created")
// 	}
// }

// func TestMemDb_Insert(t *testing.T) {
// 	memDb := peopledb.NewMemDb()

// 	person := people.NewPerson(
// 		"leorcvargas",
// 		"Léo",
// 		"1970-01-01",
// 		[]string{"golang", "nodejs"},
// 	)

// 	err := memDb.Insert(*person)
// 	if err != nil {
// 		t.Errorf("expected no errors, got %v", err)
// 	}

// 	key := person.Nickname + " " + person.Name + " " + person.StackString()

// 	got, err := memDb.Get(key)
// 	if err != nil {
// 		t.Errorf("expected no errors, got %v", err)
// 	}

// 	if got.Nickname != person.Nickname {
// 		t.Errorf("expected %s, got %s", person.Nickname, got.Nickname)
// 	}

// 	if got.Name != person.Name {
// 		t.Errorf("expected %s, got %s", person.Name, got.Name)
// 	}

// 	if got.Birthdate != person.Birthdate {
// 		t.Errorf("expected %s, got %s", person.Birthdate, got.Birthdate)
// 	}

// 	if got.StackString() != person.StackString() {
// 		t.Errorf("expected %s, got %s", person.StackString(), got.StackString())
// 	}
// }

// func TestMemDb_Get(t *testing.T) {
// 	memDb := peopledb.NewMemDb()

// 	person := people.NewPerson(
// 		"leorcvargas",
// 		"Léo",
// 		"1970-01-01",
// 		[]string{"golang", "nodejs"},
// 	)

// 	err := memDb.Insert(*person)
// 	if err != nil {
// 		t.Errorf("expected no errors, got %v", err)
// 	}

// 	key := person.Nickname + " " + person.Name + " " + person.StackString()

// 	got, err := memDb.Get(key)
// 	if err != nil {
// 		t.Errorf("expected no errors, got %v", err)
// 	}

// 	if got.Nickname != person.Nickname {
// 		t.Errorf("expected %s, got %s", person.Nickname, got.Nickname)
// 	}

// 	if got.Name != person.Name {
// 		t.Errorf("expected %s, got %s", person.Name, got.Name)
// 	}

// 	if got.Birthdate != person.Birthdate {
// 		t.Errorf("expected %s, got %s", person.Birthdate, got.Birthdate)
// 	}

// 	if got.StackString() != person.StackString() {
// 		t.Errorf("expected %s, got %s", person.StackString(), got.StackString())
// 	}
// }

// func TestMemDb_Search(t *testing.T) {
// 	list := []people.Person{
// 		{ID: uuid.New(), Nickname: "leorcvargas", Name: "Léo", Birthdate: "1970-01-01", Stack: []string{"golang", "nodejs"}},
// 		{ID: uuid.New(), Nickname: "naruto", Name: "Uzumaki Naruto", Birthdate: "1970-01-01", Stack: []string{"ruby"}},
// 		{ID: uuid.New(), Nickname: "sasuke", Name: "Uchiha Sasuke", Birthdate: "1970-01-01", Stack: []string{"python"}},
// 	}

// 	memDb := peopledb.NewMemDb()

// 	for _, person := range list {
// 		err := memDb.Insert(person)
// 		if err != nil {
// 			t.Errorf("expected no errors, got %v", err)
// 		}
// 	}

// 	testCases := []string{"rcvar", "suke", "go", "pyth", "ruby", "uzumaki", "ruto", "naruto", "sasuke", "leorcvargas"}
// 	for _, tc := range testCases {
// 		got, err := memDb.Search(tc)
// 		if err != nil {
// 			t.Errorf("expected no errors, got %v", err)
// 		}

// 		if len(got) == 0 {
// 			t.Errorf("expected at least one person for term %s, got %v", tc, got)
// 		}
// 	}
// }

func BenchmarkSearch(b *testing.B) {
	log.Println("Initializing database")

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

	log.Println("Finished initializing database")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := memDb.Search("go")
		if err != nil {
			b.Errorf("expected no errors, got %v", err)
		}
	}
}
