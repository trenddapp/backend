package migration

import (
	"embed"

	"github.com/trenddapp/backend/pkg/migration"
)

//go:embed *.sql
var fsys embed.FS

var Module = migration.NewModule(fsys)
