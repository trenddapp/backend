package migration

import (
	"context"
	"io/fs"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func CreateMigrationModule(fsys fs.FS) fx.Option {
	return fx.Invoke(
		func(
			db *bun.DB,
			logger *zap.Logger,
			migrations *migrate.Migrations,
		) error {
			if err := migrations.Discover(fsys); err != nil {
				return err
			}

			migrator := migrate.NewMigrator(db, migrations)

			if err := migrator.Init(context.TODO()); err != nil {
				return err
			}

			if _, err := migrator.Migrate(context.TODO()); err != nil {
				return err
			}

			return nil
		},
	)
}
