package service

import (
	"time"

	config_auto "github.com/vedadiyan/goal/pkg/config/auto"
	"google.golang.org/protobuf/proto"
)

func init() {
	service := New("", "", "", handler, WithCache(time.Hour))
	Register(service)
}
func handler(request proto.Message) (proto.Message, error) {
	return nil, nil
}

func init() {
	nats := config_auto.New("", true, func(value string) {

	})
	config_auto.Register(nats)
}
