package routers

import (
	"github.com/bytedance/sonic"
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
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
	})

	peopleRouter.Load(r)

	return r
}
