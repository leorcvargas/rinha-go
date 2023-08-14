package peopledb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewSearchDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./search.db")
	if err != nil {
		log.Fatal(err)
	}

	statement := `CREATE VIRTUAL TABLE people USING fts5(
		id,
		nickname,
		name,
		birthdate,
		stack,
		tokenize = "trigram case_sensitive 0"
	);`

	_, err = db.Exec(statement)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
