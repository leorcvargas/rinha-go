package peopledb

import (
	"database/sql"
	"sync/atomic"

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

func (r *roundRobinSearch) Next() TrigramSearcher {
	n := atomic.AddUint32(&r.next, 1)
	return r.searchers[(int(n)-1)%len(r.searchers)]
}

func (r *roundRobinSearch) Search(term string) ([]people.Person, error) {
	searcher := r.Next()

	return searcher.Search(term)
}

func NewRoundRobinSearcher(searchers ...TrigramSearcher) *roundRobinSearch {
	return &roundRobinSearch{
		searchers: searchers,
	}
}
