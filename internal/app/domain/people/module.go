package people

import "go.uber.org/fx"

var Module = fx.Provide(
	NewFindPeople,
	NewCreatePerson,
	NewCountPeople,
)
