package service

import (
	"sync"

	"github.com/vedadiyan/goal/pkg/runtime"
)

type ReloadStates int

const (
	RELOADING ReloadStates = iota
	RELOADED
)

var _services sync.Pool

type Service interface {
	Configure(bool)
	Start() error
	Shutdown() error
	Reload() <-chan ReloadStates
}

func Register(service Service) {
	_services.Put(service)
}

func Bootstrapper() {
	services := make([]Service, 0)
	for {
		service := _services.Get()
		if service == nil {
			break
		}
		services = append(services, service.(Service))
	}
	for _, service := range services {
		starter(service)
	}
	runtime.WaitForInterrupt(func() {
		for _, service := range services {
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
	go func(service Service) {
	LOOP:
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
						break LOOP
					}
				}
			}
		}
	}(service)
}
