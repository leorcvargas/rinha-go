package httpapi

import (
	"context"
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func NewServer(lifecycle fx.Lifecycle, router *fiber.App, _ *sql.DB) *fasthttp.Server {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Println("starting the server...")
				if err := router.Listen(":8080"); err != nil {
					log.Fatalf("error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return router.ShutdownWithContext(ctx)
		},
	})

	return router.Server()
}
