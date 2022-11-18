package bun

import (
	"crypto/tls"
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
	var tlsConfig *tls.Config

	switch cfg.SSLMode {
	case "verify-ca", "verify-full":
		tlsConfig = &tls.Config{}
	case "allow", "prefer", "require":
		tlsConfig = &tls.Config{InsecureSkipVerify: true} //nolint
	default:
		tlsConfig = nil
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(cfg.Address),
		pgdriver.WithDatabase(cfg.Database),
		pgdriver.WithPassword(cfg.Password),
		pgdriver.WithUser(cfg.User),
		pgdriver.WithTLSConfig(tlsConfig),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	if cfg.Debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db
}
