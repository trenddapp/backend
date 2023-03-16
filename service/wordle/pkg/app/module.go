package app

import (
	"go.uber.org/fx"

	"github.com/trenddapp/backend/pkg/app"
	"github.com/trenddapp/backend/pkg/database/bun"
	"github.com/trenddapp/backend/service/wordle/pkg/http"
	"github.com/trenddapp/backend/service/wordle/pkg/migration"
	"github.com/trenddapp/backend/service/wordle/pkg/repository/word"
	"github.com/trenddapp/backend/service/wordle/pkg/repository/wordle"
	"github.com/trenddapp/backend/service/wordle/pkg/workflow"
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
