package service

import (
	"log"
	"sync"
)

type ReloadStates int

const (
	RELOADING ReloadStates = iota
	RELOADED
	ACK
	ERROR
)

var _services sync.Pool

type Service interface {
	Configure(bool)
	Start() error
	Shutdown() error
	Reload() chan ReloadStates
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
	// runtime.WaitForInterrupt(func() {
	// 	for _, service := range services {
	// 		service.Shutdown()
	// 	}
	// })
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
					reloadChan <- ACK
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
