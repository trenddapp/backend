package logging

import (
	"context"

	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Driver string `yaml:"driver"`
	Level  string `yaml:"level"`
}

func NewConfig(cfg *config.YAML) (Config, error) {
	var c Config
	err := cfg.Get("logging").Populate(&c)
	return c, err
}

func NewLogger(cfg Config, lifecycle fx.Lifecycle) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	if cfg.Driver == "dev" {
		zapConfig = zap.NewDevelopmentConfig()
	}

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
