package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/config"
	"github.com/trenddapp/backend/pkg/http"
)

var BaseModule = fx.Options(
	config.Module,
	http.Module,
)
