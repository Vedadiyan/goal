package service

import (
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
	Configure() error
	Init()
	Start() error
	GetState() States
	SetState(state States)
	Stop() error
	Reload() <-chan ReloadStates
}

func Bootstrapper(services ...Service) {
	for _, service := range services {
		service.Init()
		service.Configure()
		service.Start()
		service.SetState(STARTED)
		go func(service Service) {
			for service.GetState() == STARTED {
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
								service.Configure()
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
			service.SetState(STOPPED)
			service.Stop()
		}
	})
}
