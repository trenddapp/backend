package workflow

import "go.uber.org/fx"

var BaseModule = fx.Provide(
	NewEngine,
)
