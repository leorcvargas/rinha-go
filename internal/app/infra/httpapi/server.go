package httpapi

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/config"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func startProfiling(config *config.Config) {
	log.Infof(
		"Starting CPU and Memory profiling on %s and %s",
		config.Profiling.CPU,
		config.Profiling.Mem,
	)

	cpuProfFile, err := os.Create(config.Profiling.CPU)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuProfFile)

	memoryProfFile, err := os.Create(config.Profiling.Mem)
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(memoryProfFile)

	after := time.After(3 * time.Minute)

	go func() {
		<-after
		log.Info("Stopping CPU and Memory profiling")
		pprof.StopCPUProfile()
		cpuProfFile.Close()
		memoryProfFile.Close()
	}()
}

func NewServer(
	lifecycle fx.Lifecycle,
	router *fiber.App,
	config *config.Config,
	_ *pgxpool.Pool,
) *fasthttp.Server {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				log.Info("Starting the server...")

				if config.Profiling.Enabled {
					startProfiling(config)
				}

				addr := fmt.Sprintf(":%s", config.Server.Port)
				if err := router.Listen(addr); err != nil {
					log.Fatalf("Error starting the server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping the server...")

			return router.ShutdownWithContext(ctx)
		},
	})

	return router.Server()
}
