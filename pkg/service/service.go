package service

import (
	"sync"
	"time"

	"github.com/vedadiyan/goal/pkg/runtime"
	"github.com/vedadiyan/goal/pkg/util"
)

type States int
type ReloadStates int

const (
	RUNNING      States = 1
	HEALTH_CHECK States = 1
	ERRORED      States = 2
	_STOPPED     States = 3
)

const (
	RELOADING ReloadStates = iota
	RELOADED
)

var states sync.Map

type Service interface {
	Configure(bool) error
	HealthCheck()
	Start() <-chan States
	Stop() error
	Reload() <-chan ReloadStates
}

func Bootstrapper(services ...Service) {
	for _, service := range services {
		starter(service)
	}
	runtime.WaitForInterrupt(func() {
		for _, service := range services {
			states.Store(service, _STOPPED)
			service.Stop()
		}
	})
}

func starter(service Service) {
	service.Configure(false)
	state := service.Start()
	if <-state != RUNNING {
		return
	}
	states.Store(service, RUNNING)
	go func(service Service) {
	LOOP:
		for {
			select {
			case value := <-util.GuardAgainstClosedChan(state):
				{
					switch value {
					case HEALTH_CHECK:
						{
							service.HealthCheck()
						}
					case ERRORED:
						{
							states.Store(service, ERRORED)
							break LOOP
						}
					}
				}
			case value := <-service.Reload():
				{
					switch value {
					case RELOADING:
						{
							service.Stop()
						}
					case RELOADED:
						{
							service.Configure(true)
							state = service.Start()
						}
					}
				}
			case <-time.After(time.Second):
			}
		}
	}(service)
}
