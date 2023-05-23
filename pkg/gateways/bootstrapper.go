package gateways

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var _gateways sync.Pool

func Bootstrap(app *fiber.App) {
	for {
		value := _gateways.Get()
		if value == nil {
			break
		}
		gateway := value.(func(app *fiber.App))
		gateway(app)
	}
}
