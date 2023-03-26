package config

import "go.uber.org/fx"

var BaseModule = fx.Provide(
	NewYAML,
)
