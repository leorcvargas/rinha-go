package main

import (
	"net/http"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/controllers"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/routers"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		controllers.Module,
		routers.Module,
		httpapi.Module,
		database.Module,
		people.Module,
		fx.Invoke(func(*http.Server) {}),
	)
	app.Run()
}
