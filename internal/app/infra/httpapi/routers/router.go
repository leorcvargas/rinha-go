package routers

import (
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Load()
}

func MakeRouter(
	peopleRouter *PeopleRouter,
) *fiber.App {
	r := fiber.New(fiber.Config{
		CaseSensitive: true,
		AppName:       "rinha-go by @leorcvargas",
	})

	peopleRouter.Load(r)

	return r
}
