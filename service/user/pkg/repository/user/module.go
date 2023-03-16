package user

import "go.uber.org/fx"

var BaseModule = fx.Provide(
	NewRepository,
)
