package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/app"
	"github.com/trenddapp/backend/pkg/database/bun"
	"github.com/trenddapp/backend/service/wordle/http"
	"github.com/trenddapp/backend/service/wordle/migration"
	"github.com/trenddapp/backend/service/wordle/repository/word"
	"github.com/trenddapp/backend/service/wordle/repository/wordle"
	"github.com/trenddapp/backend/service/wordle/workflow"
)

var BaseModule = fx.Options(
	app.BaseModule,
	bun.BaseModule,
	http.BaseModule,
	word.BaseModule,
	wordle.BaseModule,
	workflow.BaseModule,
)

var Module = fx.Options(
	BaseModule,
	migration.Module,
)
