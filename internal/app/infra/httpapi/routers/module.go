package routers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewPeopleRouter,
		MakeRouter,
	),
)
