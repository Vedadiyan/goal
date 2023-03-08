package service

import (
	"sync"
	"time"

	"github.com/vedadiyan/goal/pkg/runtime"
)

type States int
type ReloadStates int

const (
	STOPPED States = iota
	STARTED
)

const (
	RELOADING ReloadStates = iota
	RELOADED
)

type Service interface {
	Configure(bool) error
	Start() error
	Stop() error
	Reload() <-chan ReloadStates
}

func Bootstrapper(services ...Service) {
	var states sync.Map
	for _, service := range services {
		service.Configure(false)
		service.Start()
		states.Store(service, STARTED)
		go func(service Service) {
			for {
				value, ok := states.Load(service)
				if !ok || value.(States) == STOPPED {
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
			states.Store(service, STOPPED)
			service.Stop()
		}
	})
}
