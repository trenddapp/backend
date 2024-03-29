package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/config"
	"github.com/trenddapp/backend/pkg/http"
	"github.com/trenddapp/backend/pkg/logging"
	"github.com/trenddapp/backend/pkg/migration"
)

var BaseModule = fx.Options(
	config.BaseModule,
	http.BaseModule,
	logging.BaseModule,
	migration.BaseModule,
)
