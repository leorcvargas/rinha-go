package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/httpapi/controllers"
)

type PeopleRouter struct {
	controller *controllers.PeopleController
}

func (p *PeopleRouter) Load(r *gin.Engine) {
	r.GET("/pessoas", p.controller.Search)
	r.GET("/pessoas/:id", p.controller.Get)
	r.GET("/contagem-pessoas", p.controller.CountAll)
	r.POST("/pessoas", p.controller.Create)
}

func NewPeopleRouter(
	controller *controllers.PeopleController,
) *PeopleRouter {
	return &PeopleRouter{
		controller: controller,
	}
}
