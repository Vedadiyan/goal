package service

import (
	"log"
	"sync"

	"github.com/vedadiyan/goal/pkg/runtime"
)

type ReloadStates int

const (
	RELOADING ReloadStates = iota
	RELOADED
	READY
	ERROR
)

var _services []any
var _mute sync.Mutex
var _skipInterrupt bool

type Service interface {
	Configure(bool)
	Start() error
	Shutdown() error
	Reload() chan ReloadStates
}

func init() {
	_mute.Lock()
	_services = make([]any, 0)
	_mute.Unlock()
}

func Register(service Service) {
	_services = append(_services, service)
}

func Bootstrap() {
	services := make([]Service, 0)
	for _, service := range _services {
		services = append(services, service.(Service))
	}
	for _, service := range services {
		starter(service)
	}
	if !_skipInterrupt {
		runtime.WaitForInterrupt(func() {
			for _, service := range services {
				_ = service.Shutdown()

			}
		})
	}
}

func starter(service Service) {
	log.Println("configuring")
	service.Configure(false)
	log.Println("configured")
	log.Println("starting")
	err := service.Start()
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("started")
	go func(service Service) {
		reloadChan := service.Reload()
	LOOP:
		for value := range reloadChan {
			switch value {
			case RELOADING:
				{
					log.Println("reloading")
					err := service.Shutdown()
					if err != nil {
						reloadChan <- ERROR
						log.Fatalln(err)
						return
					}
					reloadChan <- READY
					log.Println("reloading done")
				}
			case RELOADED:
				{
					log.Println("reconfiguring")
					service.Configure(true)
					log.Println("reconfigured")
					log.Println("restarting")
					err := service.Start()
					if err != nil {
						log.Fatalln(err)
						break LOOP
					}
					log.Println("restarted")
				}
			}
		}
	}(service)
}
