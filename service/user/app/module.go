package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/app"
	"github.com/trenddapp/backend/pkg/database/bun"
	"github.com/trenddapp/backend/service/user/http"
	"github.com/trenddapp/backend/service/user/migration"
	"github.com/trenddapp/backend/service/user/repository/nonce"
	"github.com/trenddapp/backend/service/user/repository/user"
)

var BaseModule = fx.Options(
	app.BaseModule,
	bun.BaseModule,
	http.BaseModule,
	nonce.BaseModule,
	user.BaseModule,
)

var Module = fx.Options(
	BaseModule,
	migration.Module,
)
