package gateways

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vedadiyan/goal/pkg/protoutil"
	"github.com/vedadiyan/goal/pkg/proxy"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func GetJSONReq[T proto.Message](c *fiber.Ctx, req T) error {
	values := make(map[string]any)
	for _, key := range c.Route().Params {
		values[key] = c.Params(key)
	}
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		values[string(key)] = string(value)
	})
	if len(c.Body()) != 0 {
		c.BodyParser(&values)
	}
	out, err := json.Marshal(values)
	if err != nil {
		return err
	}
	err = protojson.Unmarshal(out, req)
	if err != nil {
		return err
	}
	return nil
}

func Single[TRequest any, TResponse any](app *fiber.App, uri string, method string, to string) *proxy.NATSProxy[proto.Message] {
	proxy := proxy.Create[TResponse]("$etcd:/nats", to)
	app.Add(method, uri, func(c *fiber.Ctx) error {
		var inst TRequest
		req := any(&inst).(proto.Message)
		err := GetJSONReq(c, req)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		res, err := proxy.Send(req)
		if err != nil {
			return err
		}
		mapper, err := protoutil.Marshal(any(*res).(proto.Message))
		if err != nil {
			return nil
		}
		return c.JSON(mapper)
	})
	return proxy
}

func Aggregated[TRequest any, TResponse any](app *fiber.App, uri string, method string, to map[string]string) map[string]*proxy.NATSProxy[proto.Message] {
	proxies := make(map[string]*proxy.NATSProxy[proto.Message])
	for key, gateway := range to {
		correctedURI := strings.TrimSuffix(uri, "/")
		proxies[key] = Single[TRequest, TResponse](app, fmt.Sprintf("%s/%s", correctedURI, key), method, gateway)
	}
	app.Add(method, uri, func(c *fiber.Ctx) error {
		var inst TRequest
		req := any(&inst).(proto.Message)
		err := GetJSONReq(c, req)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		out := make(map[string]map[string]any)
		for key, proxy := range proxies {
			res, err := proxy.Send(req)
			if err != nil {
				out[key] = map[string]any{
					"error": err.Error(),
				}
				continue
			}
			mapper, err := protoutil.Marshal(any(*res).(proto.Message))
			if err != nil {
				out[key] = map[string]any{
					"error": err.Error(),
				}
				continue
			}
			out[key] = mapper
		}
		return c.JSON(out)
	})
	return proxies
}

func Forward[TRequest any, TResponse any](uri string, method string, to any) {
	_gateways.Put(func(app *fiber.App) {
		switch t := to.(type) {
		case string:
			{
				Single[TRequest, TResponse](app, uri, method, t)
			}
		case map[string]string:
			{
				Aggregated[TRequest, TResponse](app, uri, method, t)
			}
		}
	})
}
