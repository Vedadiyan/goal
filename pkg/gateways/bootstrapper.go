package gateways

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var _gateways sync.Pool

func Register(gateway *FiberGateway) {
	_gateways.Put(gateway)
}
func Bootstrap(app *fiber.App) {
	for {
		value := _gateways.Get()
		if value == nil {
			break
		}
		gateway := value.(*FiberGateway)
		gateway.initializer()
		for key, endpoint := range gateway.endpoints {
			for method, handler := range endpoint {
				app.Add(method, fmt.Sprintf("%s/%s", gateway.route, key), handler)
			}
		}
	}
}
