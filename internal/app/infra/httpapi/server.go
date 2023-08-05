package httpapi

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewServer(lifecycle fx.Lifecycle, router *echo.Echo, _ *gorm.DB) *http.Server {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go router.Start(":8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return router.Shutdown(ctx)
		},
	})

	return router.Server
}
