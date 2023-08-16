package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/controllers"
)

type PeopleRouter struct {
	controller *controllers.PeopleController
}

func (p *PeopleRouter) Load(r *fiber.App) {
	r.Get("/pessoas", p.controller.Search)
	r.Get("/pessoas/:id", p.controller.Get)
	r.Get("/contagem-pessoas", p.controller.CountAll)
	r.Post("/pessoas", p.controller.Create)
}

func NewPeopleRouter(
	controller *controllers.PeopleController,
) *PeopleRouter {
	return &PeopleRouter{
		controller: controller,
	}
}
