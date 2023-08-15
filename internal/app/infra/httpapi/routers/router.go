package routers

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Load()
}

func MakeRouter(
	peopleRouter *PeopleRouter,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	peopleRouter.Load(r)

	return r
}
