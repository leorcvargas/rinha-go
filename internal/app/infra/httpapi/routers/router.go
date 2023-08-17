package routers

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
		Prefork:       true,
	}

	if os.Getenv("ENABLE_SONIC_JSON") == "1" {
		log.Info("Loading Sonic JSON into the router")
		cfg.JSONEncoder = sonic.Marshal
		cfg.JSONDecoder = sonic.Unmarshal
	}

	r := fiber.New()

	peopleRouter.Load(r)

	return r
}
