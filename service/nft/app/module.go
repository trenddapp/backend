package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/app"
	"github.com/trenddapp/backend/service/nft/client/etherscan"
	"github.com/trenddapp/backend/service/nft/client/nftport"
	"github.com/trenddapp/backend/service/nft/http"
)

var BaseModule = fx.Options(
	app.BaseModule,
	etherscan.Module,
	http.Module,
	nftport.Module,
)
