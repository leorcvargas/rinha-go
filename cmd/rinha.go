package main

import (
	"log"
	"net/http"

	"github.com/leorcvargas/rinha-2023-q3/internal/app/domain/people"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/database"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/controllers"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/routers"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/pubsub"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/worker"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		controllers.Module,
		routers.Module,
		httpapi.Module,
		database.Module,
		people.Module,
		worker.Module,
		pubsub.Module,
		fx.Invoke(func(worker *worker.Inserter) {
			log.Println("Starting worker.Inserter")
			go worker.Run()
		}),
		fx.Invoke(func(subscriber *pubsub.PersonInsertSubscriber) {
			log.Println("Starting pubsub.Subscriber")
			go subscriber.Subscribe()
		}),
		fx.Invoke(func(*http.Server) {}),
	)
	app.Run()
}
