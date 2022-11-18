package migration

import (
	"github.com/uptrace/bun/migrate"
	"go.uber.org/fx"
)

var BaseModule = fx.Provide(
	migrate.NewMigrations,
)
