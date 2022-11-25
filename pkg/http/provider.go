package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/config"
	"go.uber.org/fx"
)

type Config struct {
	Port int `yaml:"port"`
}

func NewConfig(cfg *config.YAML) (*Config, error) {
	c := &Config{}
	if err := cfg.Get("http").Populate(c); err != nil {
		return nil, err
	}

	return c, nil
}

func NewRouter(cfg *Config, lifecycle fx.Lifecycle) *gin.Engine {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsHandler := cors.New(corsConfig)

	router := gin.Default()
	router.Use(corsHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil {
					// TODO: Log error.
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			return server.Close()
		},
	})

	return router
}
