package main

import (
	"github.com/google/uuid"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database/peopledb"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/controllers"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/routers"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"

	_ "go.uber.org/automaxprocs"
)

func main() {
	uuid.EnableRandPool()

	app := fx.New(
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
