package logging

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var BaseModule = fx.Options(
	fx.Provide(ProvideConfig),
	fx.Provide(ProvideZapConfig),
	fx.Provide(ProvideLogger),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)
