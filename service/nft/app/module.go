package app

import (
	"go.uber.org/fx"

	"github.com/dapp-z/backend/pkg/app"
	"github.com/dapp-z/backend/service/nft/client/etherscan"
	"github.com/dapp-z/backend/service/nft/client/nftport"
	"github.com/dapp-z/backend/service/nft/http"
)

var BaseModule = fx.Options(
	app.BaseModule,
	etherscan.Module,
	http.Module,
	nftport.Module,
)
