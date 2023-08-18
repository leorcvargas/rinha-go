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
	cfg := fiber.Config{
		AppName:       "rinha-go by @leorcvargas",
		CaseSensitive: true,
		JSONEncoder:   sonic.Marshal,
		JSONDecoder:   sonic.Unmarshal,
	}

	// if os.Getenv("ENABLE_SONIC_JSON") == "1" {
	// 	log.Info("Loading Sonic JSON into the router")
	// 	cfg.JSONEncoder = sonic.Marshal
	// 	cfg.JSONDecoder = sonic.Unmarshal
	// }

	r := fiber.New(cfg)

	peopleRouter.Load(r)

	return r
}
