package worker

import "go.uber.org/fx"

var Module = fx.Provide(
	NewInserter,
	CreateInsertChannel,
)
