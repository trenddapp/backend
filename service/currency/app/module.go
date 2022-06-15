package app

import (
	"go.uber.org/fx"

	"github.com/dapp-z/backend/pkg/app"
	"github.com/dapp-z/backend/service/currency/client/coinmarketcap"
	"github.com/dapp-z/backend/service/currency/http"
)

var BaseModule = fx.Options(
	app.BaseModule,
	coinmarketcap.Module,
	http.Module,
)
