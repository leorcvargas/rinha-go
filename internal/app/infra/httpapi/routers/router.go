package routers

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/leorcvargas/rinha-2023-q3/internal/app/infra/config"
)

type Router interface {
	Load()
}

func MakeRouter(
	peopleRouter *PeopleRouter,
	config *config.Config,
) *fiber.App {
	cfg := fiber.Config{
		AppName:       "rinha-go by @leorcvargas",
		CaseSensitive: true,
	}

	if config.Server.UseSonic {
		log.Info("Loading Sonic JSON into the router")
		cfg.JSONEncoder = sonic.Marshal
		cfg.JSONDecoder = sonic.Unmarshal
	}

	r := fiber.New(cfg)

	peopleRouter.Load(r)

	return r
}
