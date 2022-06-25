package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/app"
	"github.com/trenddapp/backend/service/currency/client/coinmarketcap"
	"github.com/trenddapp/backend/service/currency/http"
)

var BaseModule = fx.Options(
	app.BaseModule,
	coinmarketcap.Module,
	http.Module,
)
