package peopledb

import (
	"database/sql"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
)

type TrigramSearcher interface {
	Search(term string) ([]people.Person, error)
}

type LocalTrigramSearcher struct {
	db *sql.DB
}

func (l *LocalTrigramSearcher) Search(term string) ([]people.Person, error) {
	rows, err := l.db.Query(
		`SELECT * FROM people(?) LIMIT 50;`,
		term,
	)
	if err != nil {
		return nil, err
	}

	result, err := mapSearchResult(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type PsqlTrigramSearcher struct {
	db *sql.DB
}

func (p *PsqlTrigramSearcher) Search(term string) ([]people.Person, error) {
	rows, err := p.db.Query(
		SearchPeopleTrgmQuery,
		term,
	)
	if err != nil {
		return nil, err
	}

	result, err := mapSearchResult(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type roundRobinSearch struct {
	searchers []TrigramSearcher
	next      uint32
}

func (r *roundRobinSearch) Search(term string) ([]people.Person, error) {
	searcher := r.searchers[r.next%uint32(len(r.searchers))]
	r.next++

	return searcher.Search(term)
}

func NewRoundRobinSearcher(searchers ...TrigramSearcher) *roundRobinSearch {
	return &roundRobinSearch{
		searchers: searchers,
	}
}
