package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/trenddapp/backend/service/currency/pkg/app"
	"github.com/trenddapp/backend/service/currency/pkg/http"
)

var serveHTTPCmd = &cobra.Command{
	Use: "http",
	Run: func(cmd *cobra.Command, args []string) {
		fx.New(app.BaseModule, fx.Invoke(http.RegisterRoutes)).Run()
	},
}

func init() {
	serveCmd.AddCommand(serveHTTPCmd)
}
