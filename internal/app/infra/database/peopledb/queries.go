package peopledb

const (
	InsertPersonQuery = "INSERT INTO people (id, nickname, name, birthdate, stack, search) VALUES ($1, $2, $3, $4, $5, $6);"

	SelectPersonByIDQuery = "SELECT id, nickname, name, birthdate, stack FROM people WHERE id = $1;"

	SearchPeopleTrgmQuery = `SELECT id, nickname, name, birthdate, stack FROM people p
	WHERE p.search LIKE '%' || $1 || '%'
	LIMIT 50;`

	CountPeopleQuery = "SELECT COUNT(*) FROM people;"
)
