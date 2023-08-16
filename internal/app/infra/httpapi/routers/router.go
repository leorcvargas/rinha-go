package routers

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	Load()
}

func MakeRouter(
	peopleRouter *PeopleRouter,
) *fiber.App {
	cfg := fiber.Config{
		AppName:       "rinha-go by @leorcvargas",
		CaseSensitive: true,
	}

	if os.Getenv("ENABLE_SONIC_JSON") == "1" {
		cfg.JSONEncoder = sonic.Marshal
		cfg.JSONDecoder = sonic.Unmarshal
	}

	r := fiber.New()

	peopleRouter.Load(r)

	return r
}
