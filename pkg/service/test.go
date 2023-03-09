package service

import (
	"github.com/nats-io/nats.go"
	"github.com/vedadiyan/goal/pkg/di"
)

type TestService struct {
	conn         *nats.Conn
	reloadState  chan ReloadStates
	subscription *nats.Subscription
}

func (t *TestService) Configure(b bool) {
	if !b {
		di.OnSingletonRefreshWithName("nats", func(e di.Events) {
			if e == di.REFRESHED {
				t.conn = *di.ResolveWithNameOrPanic[*nats.Conn]("nats", nil)
				t.reloadState <- RELOADED
				return
			}
			t.reloadState <- RELOADING
		})
		return
	}
	t.conn = *di.ResolveWithNameOrPanic[*nats.Conn]("nats", nil)
}
func (t *TestService) Start() error {
	subs, err := t.conn.QueueSubscribe("abcd.*", "balanced", func(msg *nats.Msg) {
		go func() {
			switch msg.Subject {
			case "abcd.health_check":
				{

				}
			case "abcd.service":
				{

				}
			}
		}()
	})
	if err != nil {
		return err
	}
	t.subscription = subs
	return nil
}
func (t TestService) Shutdown() error {
	return t.subscription.Unsubscribe()
}
func (t TestService) Reload() <-chan ReloadStates {
	return t.reloadState
}
