package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router interface {
	Load()
}

func MakeRouter(
	peopleRouter *PeopleRouter,
) *echo.Echo {
	router := echo.New()

	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.Secure())
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	peopleRouter.Load(router)

	return router
}
