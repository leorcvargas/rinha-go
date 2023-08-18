package peopledb

const (
	InsertPersonQuery = "INSERT INTO people (id, nickname, name, birthdate, stack, trgm_q) VALUES ($1, $2, $3, $4, $5, $6);"

	SelectPersonByIDQuery = "SELECT id, nickname, name, birthdate, stack FROM people WHERE id = $1;"

	SearchPeopleFtsQuery = `SELECT id, nickname, name, birthdate, stack FROM people p
	WHERE p.fts_q @@ plainto_tsquery('people_terms', $1)
	LIMIT 50;`

	SearchPeopleTrgmQuery = `SELECT id, nickname, name, birthdate, stack FROM people p
	WHERE p.trgm_q ILIKE '%' || $1 || '%'
	LIMIT 50;`

	CountPeopleQuery = "SELECT COUNT(*) FROM people;"
)
