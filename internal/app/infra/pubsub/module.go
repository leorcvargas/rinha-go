package pubsub

import "go.uber.org/fx"

var Module = fx.Provide(
	NewPersonInsertSubscriber,
)
