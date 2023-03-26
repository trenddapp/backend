package http

import "go.uber.org/fx"

var BaseModule = fx.Provide(
	NewConfig,
	NewRouter,
)
