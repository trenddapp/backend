package etherscan

import "go.uber.org/fx"

var Module = fx.Provide(
	NewClient,
	NewConfig,
)
