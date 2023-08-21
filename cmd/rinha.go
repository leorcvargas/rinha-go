package main

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/config"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/controllers"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/routers"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"

	_ "net/http/pprof"

	_ "go.uber.org/automaxprocs"
)

func main() {
	uuid.EnableRandPool()

	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("Coudn't load .env file")
	}

	app := fx.New(
		config.Module,
		controllers.Module,
		routers.Module,
		httpapi.Module,
		database.Module,
		people.Module,
		fx.Invoke(func(dispatcher *peopledb.Dispatcher) {
			go dispatcher.Run()
		}),
		fx.Invoke(func(*fasthttp.Server) {}),
		fx.NopLogger,
	)

	app.Run()
}
