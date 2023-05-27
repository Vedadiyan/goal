package gateways

import (
	"github.com/gofiber/fiber/v2"
)

func Register(uri string, method string, handler func(ctx *fiber.Ctx) error) {
	_gateways = append(_gateways, func(app *fiber.App) {
		app.Add(method, uri, handler)
	})
}
