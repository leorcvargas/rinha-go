package peopledb

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewLocalDatabase() *sql.DB {
	os.Remove("./localsearch.db")
	os.Create("./localsearch.db")

	db, err := sql.Open("sqlite3", "./localsearch.db")
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
