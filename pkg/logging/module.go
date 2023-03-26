package logging

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var BaseModule = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewLogger),
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger}
	}),
)
