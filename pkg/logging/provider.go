package logging

import (
	"context"

	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ProvideConfig(cfg *config.YAML) (Config, error) {
	var c Config
	err := cfg.Get("logging").Populate(&c)
	return c, err
}

func ProvideZapConfig(cfg Config) zap.Config {
	if cfg.Driver == "dev" {
		return zap.NewDevelopmentConfig()
	}

	return zap.NewProductionConfig()
}

func ProvideLogger(cfg Config, lifecycle fx.Lifecycle, zapConfig zap.Config) (*zap.Logger, error) {
	if cfg.Level != "" {
		var level zapcore.Level

		if err := level.Set(cfg.Level); err != nil {
			return nil, err
		}

		zapConfig.Level.SetLevel(level)
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			// https://github.com/uber-go/zap/issues/370
			logger.Sync() //nolint:errcheck
			return nil
		},
	})

	return logger, nil
}
