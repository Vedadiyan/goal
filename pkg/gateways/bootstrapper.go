package gateways

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	_gateways []func(app *fiber.App)
	_mut      sync.Mutex
)

func Bootstrap(app *fiber.App) {
	for _, value := range _gateways {
		value(app)
	}
}
