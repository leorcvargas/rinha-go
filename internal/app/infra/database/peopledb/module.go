package peopledb

import (
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewPeopleDbCache,
	NewMemDb,
	NewMem2,
	fx.Annotate(
		NewPersonRepository,
		fx.As(new(people.Repository)),
	),
)
