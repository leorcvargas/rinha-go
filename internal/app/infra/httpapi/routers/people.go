package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/controllers"
)

type PeopleRouter struct {
	controller *controllers.PeopleController
}

func (p *PeopleRouter) Load(parent *echo.Echo) {
	parent.GET("/pessoas", p.controller.GetAll)
	parent.POST("/pessoas", p.controller.Create)
}

func NewPeopleRouter(
	controller *controllers.PeopleController,
) *PeopleRouter {
	return &PeopleRouter{
		controller: controller,
	}
}
