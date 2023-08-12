package httpapi

import (
	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewServer,
)
