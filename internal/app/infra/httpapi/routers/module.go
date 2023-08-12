package routers

import "go.uber.org/fx"

var Module = fx.Provide(
	NewPeopleRouter,
	MakeRouter,
)
