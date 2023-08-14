package peopledb

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewLocalDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		log.Fatal(err)
	}

	statement := `CREATE VIRTUAL TABLE IF NOT EXISTS people USING fts5(
		id UNINDEXED,
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
