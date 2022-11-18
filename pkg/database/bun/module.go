package bun

import (
	"go.uber.org/fx"
)

var BaseModule = fx.Options(
	fx.Provide(ProvideConfig),
	fx.Provide(ProvideDB),
)
