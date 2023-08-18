package httpapi

import (
	"context"
	"database/sql"
	"os"
	"runtime/pprof"

	"github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func NewServer(lifecycle fx.Lifecycle, router *fiber.App, _ *sql.DB) *fasthttp.Server {
	// profile cpu
	f, err := os.Create(os.Getenv("CPU_PROFILE"))
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// profile mem
	mf, err := os.Create(os.Getenv("MEM_PROFILE"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.WriteHeapProfile(mf)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info("Starting the server...")
				if err := router.Listen(":8080"); err != nil {
					log.Fatalf("Error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer func() {
				pprof.StopCPUProfile()
				f.Close()
				mf.Close()
			}()

			log.Info("Stopping the server...")

			return router.ShutdownWithContext(ctx)
		},
	})

	return router.Server()
}
