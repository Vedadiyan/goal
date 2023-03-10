package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vedadiyan/goal/pkg/di"
	"google.golang.org/protobuf/proto"
)

func TestService(t *testing.T) {
	_skipInterrupt = true
	di.AddSinletonWithName("nats", func() (*nats.Conn, error) {
		return nats.Connect("127.0.0.1:4222")
	})
	service := New("nats", "test", "test", handler, WithCache(time.Hour))
	service2 := New("nats", "test2", "test", handler, WithCache(time.Hour))
	Register(service)
	Register(service2)
	Bootstrapper()
	<-time.After(time.Second)
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		}, func(conn *nats.Conn) {
			conn.Close()
		})
	}()
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		}, func(conn *nats.Conn) {
			conn.Close()
		})
	}()
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		}, func(conn *nats.Conn) {
			conn.Close()
		})
	}()
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		}, func(conn *nats.Conn) {
			conn.Close()
		})
	}()
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		}, func(conn *nats.Conn) {
			conn.Close()
		})
	}()
	go func() {
		di.RefreshSinletonWithName("nats", func() (*nats.Conn, error) {
			return nats.Connect("127.0.0.1:4222")
		}, func(conn *nats.Conn) {
			conn.Close()
		})
	}()
	<-time.After(time.Second * 10)
}

func handler(request proto.Message) (proto.Message, error) {
	return nil, nil
}

func TestMap(t *testing.T) {
	i := 0
	x := &i
	di.AddSinleton(func() (instance *int, err error) {
		return x, nil
	})
	go func() {
		value := di.ResolveOrPanic[*int](nil)
		for {
			fmt.Println(**value)
			<-time.After(time.Second)
		}
	}()
	<-time.After(time.Second * 5)
	i = 10
	<-time.After(time.Second * 5)
}
