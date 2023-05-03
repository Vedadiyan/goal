package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/vedadiyan/goal/pkg/di"
)

// func TestService(t *testing.T) {
// 	_skipInterrupt = true
// 	di.AddSinletonWithName("nats", func() (*nats.Conn, error) {
// 		return nats.Connect("127.0.0.1:4222")
// 	})
// 	for i := 0; i < 100; i++ {
// 		service := New("nats", fmt.Sprintf("%d", i), "test", handler, WithCache(time.Hour))
// 		Register(service)
// 	}
// 	Bootstrapper()
// 	<-time.After(time.Second)
// 	go func() {
// 		_, err := di.RefreshSinletonWithName("nats", func(current *nats.Conn) (*nats.Conn, error) {
// 			current.Drain()
// 			return nats.Connect("127.0.0.1:4222")
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
// 	<-time.After(time.Second * 10)
// }

// func handler(request proto.Message) (proto.Message, error) {
// 	return nil, nil
// }

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
