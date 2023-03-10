package gateways

import (
	"github.com/gofiber/fiber/v2"
)

type FiberGateway struct {
	route     string
	endpoints map[string]map[string]func(c *fiber.Ctx) error
}

func (a *FiberGateway) Add(name string, method string, fn func(c *fiber.Ctx) error) {
	if _, ok := a.endpoints[name]; !ok {
		a.endpoints[name] = make(map[string]func(c *fiber.Ctx) error)
	}
	a.endpoints[name][method] = fn
}

func New(route string) *FiberGateway {
	apiGateway := FiberGateway{
		route:     route,
		endpoints: make(map[string]map[string]func(c *fiber.Ctx) error),
	}
	return &apiGateway
}
