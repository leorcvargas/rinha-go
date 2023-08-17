package peopledb

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func top(label string) func() {
	x := time.Now()

	return func() {
		if time.Since(x) > 100*time.Millisecond {
			log.Warn(label, " time: ", time.Since(x))
		}
	}
}
