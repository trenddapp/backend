package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/app"
	"github.com/trenddapp/backend/service/currency/pkg/client/coinmarketcap"
	"github.com/trenddapp/backend/service/currency/pkg/http"
)

var BaseModule = fx.Options(
	app.BaseModule,
	coinmarketcap.Module,
	http.Module,
)
