package gateways

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/vedadiyan/goal/pkg/protoutil"
	protoval "github.com/vedadiyan/goal/pkg/protoval"
	"github.com/vedadiyan/goal/pkg/proxy"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Gateway struct {
	useMeta         bool
	validationField *string
}

type GatewayOption func(gateway *Gateway)

func GetJSONReq[T proto.Message](c *fiber.Ctx, req T, useMeta bool) error {
	values := make(map[string]any)
	if len(c.Body()) != 0 {
		err := c.BodyParser(&values)
		if err != nil {
			return err
		}
	}
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		values[string(key)] = string(value)
	})
	for _, key := range c.Route().Params {
		values[key] = c.Params(key)
	}
	if useMeta {
		meta := make(map[string]any)
		meta["remote_ip"] = c.Context().RemoteIP()
		values["meta"] = meta
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

func Single[TRequest any, TResponse any](app *fiber.App, uri string, method string, to string, options ...GatewayOption) *proxy.NATSProxy[proto.Message] {
	proxy := proxy.Create[TResponse]("default_nats", to)
	var gateway Gateway
	for _, option := range options {
		option(&gateway)
	}
	app.Add(method, uri, func(c *fiber.Ctx) error {
		var inst TRequest
		req := any(&inst).(proto.Message)
		err := GetJSONReq(c, req, gateway.useMeta)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		if gateway.validationField != nil {
			vc := protoval.New(*gateway.validationField, req)
			err := vc.Validate()
			if err != nil {
				return err
			}
			if !vc.IsValid() {
				c.Status(400)
				data := vc.Errors()
				return c.JSON(map[string]any{
					"request": req,
					"errors":  data,
				})
			}
		}
		res, err := proxy.Send(req)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			c.Response().Header.Add("Content-Type", "application/json")
			c.Send([]byte(err.Error()))
			return nil
		}
		mapper, err := protoutil.Marshal(any(*res).(proto.Message))
		if err != nil {
			return nil
		}
		return c.JSON(mapper)
	})
	return proxy
}

func Aggregated[TRequest any, TResponse any](app *fiber.App, uri string, method string, to map[string]string, options ...GatewayOption) map[string]*proxy.NATSProxy[proto.Message] {
	proxies := make(map[string]*proxy.NATSProxy[proto.Message])
	for key, gateway := range to {
		correctedURI := strings.TrimSuffix(uri, "/")
		proxies[key] = Single[TRequest, TResponse](app, fmt.Sprintf("%s/%s", correctedURI, key), method, gateway, options...)
	}
	var gateway Gateway
	for _, option := range options {
		option(&gateway)
	}
	app.Add(method, uri, func(c *fiber.Ctx) error {
		out := make(map[string]map[string]any)
		var inst TRequest
		req := any(&inst).(proto.Message)
		err := GetJSONReq(c, req, gateway.useMeta)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		if gateway.validationField != nil {
			vc := protoval.New(*gateway.validationField, req)
			err := vc.Validate()
			if err != nil {
				return err
			}
			if !vc.IsValid() {
				c.Status(400)
				return c.JSON(vc.Errors())
			}
		}
		var wg sync.WaitGroup
		var mut sync.Mutex
		for key, proxy_ := range proxies {
			wg.Add(1)
			go func(key string, proxy *proxy.NATSProxy[protoreflect.ProtoMessage]) {
				defer mut.Unlock()
				defer wg.Done()
				res, err := proxy.Send(req)
				if err != nil {
					mut.Lock()
					out[key] = map[string]any{
						"error": err.Error(),
					}
					return
				}
				mapper, err := protoutil.Marshal(any(*res).(proto.Message))
				if err != nil {
					mut.Lock()
					out[key] = map[string]any{
						"error": err.Error(),
					}
					return
				}
				mut.Lock()
				out[key] = mapper
			}(key, proxy_)
		}
		wg.Wait()
		return c.JSON(out)
	})
	return proxies
}

func Forward[TRequest any, TResponse any](uri string, method string, to any, options ...GatewayOption) {
	_gateways = append(_gateways, func(app *fiber.App) {
		switch t := to.(type) {
		case string:
			{
				Single[TRequest, TResponse](app, uri, method, t, options...)
			}
		case map[string]string:
			{
				Aggregated[TRequest, TResponse](app, uri, method, t, options...)
			}
		}
	})
}

func UseMeta() GatewayOption {
	return func(gateway *Gateway) {
		gateway.useMeta = true
	}
}

func UseValidation(validationField string) GatewayOption {
	return func(gateway *Gateway) {
		gateway.validationField = &validationField
	}
}
