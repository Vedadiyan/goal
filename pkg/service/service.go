package service

import (
	"sync"
	"time"

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

type Service interface {
	Configure(bool) error
	HealthCheck()
	Start() <-chan States
	Stop() error
	Reload() <-chan ReloadStates
}

func Bootstrapper(services ...Service) {
	var states sync.Map
	for _, service := range services {
		service.Configure(false)
		state := service.Start()
		if <-state != RUNNING {
			continue
		}
		states.Store(service, RUNNING)
		go func(service Service) {
		LOOP:
			for value := range state {
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
		}(service)
		go func(service Service) {
			for {
				value, ok := states.Load(service)
				if !ok || value.(States) > RUNNING {
					break
				}
				select {
				case reloadState := <-service.Reload():
					{
						switch reloadState {
						case RELOADING:
							{
								service.Stop()
							}
						case RELOADED:
							{
								service.Configure(true)
								service.Start()
							}
						}
					}
				case <-time.After(time.Second):
				}
			}
		}(service)
	}
	runtime.WaitForInterrupt(func() {
		for _, service := range services {
			states.Store(service, _STOPPED)
			service.Stop()
		}
	})
}
