package bun

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"go.uber.org/config"
	"go.uber.org/zap"
)

func ProvideConfig(cfg *config.YAML) (Config, error) {
	var c Config
	err := cfg.Get("databases.pg").Populate(&c)
	return c, err
}

func ProvideDB(cfg Config, logger *zap.Logger) *bun.DB {
	db := bun.NewDB(
		sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN))),
		pgdialect.New(),
	)

	if cfg.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db
}
