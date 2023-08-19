package httpapi

import (
	"context"
	"os"
	"runtime/pprof"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func prof() func() {
	f, err := os.Create(os.Getenv("CPU_PROFILE"))
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)

	mf, err := os.Create(os.Getenv("MEM_PROFILE"))
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(mf)

	return func() {
		pprof.StopCPUProfile()
		f.Close()
		mf.Close()
	}
}

func NewServer(lifecycle fx.Lifecycle, router *fiber.App, _ *pgxpool.Pool) *fasthttp.Server {
	var shutdown func()

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				shutdown = prof()
				log.Info("Starting the server...")
				if err := router.Listen(":8080"); err != nil {
					log.Fatalf("Error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer shutdown()
			log.Info("Stopping the server...")

			return router.ShutdownWithContext(ctx)
		},
	})

	return router.Server()
}
