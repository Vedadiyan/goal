package service

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vedadiyan/goal/pkg/di"
	"google.golang.org/protobuf/proto"
)

func TestService(t *testing.T) {
	di.AddSinletonWithName("nats", func() (*nats.Conn, error) {
		return nats.Connect("127.0.0.1:4222")
	})
	service := New("nats", "test", "test", handler, WithCache(time.Hour))
	Register(service)
	Bootstrapper()
	<-time.After(time.Second)
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		})
	}()
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		})
	}()
	<-time.After(time.Hour)
}

func handler(request proto.Message) (proto.Message, error) {
	return nil, nil
}
