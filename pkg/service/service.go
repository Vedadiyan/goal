package service

import (
	"sync"

	"github.com/vedadiyan/goal/pkg/runtime"
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
	Configure(bool)
	Start() error
	Shutdown() error
	Reload() <-chan ReloadStates
}

func Bootstrapper(services ...Service) {
	for _, service := range services {
		starter(service)
	}
	runtime.WaitForInterrupt(func() {
		for _, service := range services {
			states.Store(service, _STOPPED)
			service.Shutdown()
		}
	})
}

func starter(service Service) {
	service.Configure(false)
	err := service.Start()
	if err != nil {
		return
	}
	states.Store(service, RUNNING)
	go func(service Service) {
		for value := range service.Reload() {
			switch value {
			case RELOADING:
				{
					service.Shutdown()
				}
			case RELOADED:
				{
					service.Configure(true)
					err := service.Start()
					if err != nil {
						states.Store(service, ERRORED)
					}
				}
			}
		}
	}(service)
}
