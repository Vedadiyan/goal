package config_auto

import (
	"context"
	"sync"

	"github.com/vedadiyan/goal/pkg/config"
	config_etcd "github.com/vedadiyan/goal/pkg/config/etcd"
)

type Initializer interface {
	Init(value any) error
}

type String struct {
	Key   string
	Watch bool
	CB    func(value string)
}

type KeyValue map[string]any

type Object struct {
	Key   string
	Watch bool
	CB    func(value KeyValue)
}

var _initializers sync.Pool

func (k KeyValue) GetStringValue(key string) (string, error) {
	if value, ok := k[key]; ok {
		if str, ok := value.(string); ok {
			return str, nil
		}
		return "", config.INVALID_OBJECT
	}
	return "", config.KEY_NOT_FOUND
}

func (s String) Init(value any) error {
	if str, ok := value.(string); ok {
		s.CB(str)
		return nil
	}
	return config.INVALID_OBJECT
}
func (o Object) Init(value any) error {
	if keyValue, ok := value.(map[string]any); ok {
		o.CB(keyValue)
		return nil
	}
	return config.INVALID_OBJECT
}

func Register(initializer Initializer) {
	_initializers.Put(initializer)
}

func Bootstrap(url string) error {
	etcdCnfxReader, err := config_etcd.NewClient([]string{url})
	if err != nil {
		return err
	}
	for {
		value := _initializers.Get()
		if value == nil {
			break
		}
		initializer := value.(Initializer)
		switch t := initializer.(type) {
		case String:
			{
				etcdCnfx, err := etcdCnfxReader.ReadKey(context.TODO(), t.Key)
				if err != nil {
					return err
				}
				value, err := etcdCnfx.GetString(t.Key)
				if err != nil {
					return err
				}
				if t.CB != nil {
					err := t.Init(value)
					if err != nil {
						return err
					}
					if t.Watch {
						etcdCnfxReader.Watch(context.TODO(), t.Key, func(etcdValue *config_etcd.EtcdValue, err error) {
							if err != nil {
								return
							}
							t.Init(value)
						})
					}
				}
			}
		case Object:
			{
				etcdCnfx, err := etcdCnfxReader.ReadKey(context.TODO(), t.Key)
				if err != nil {
					return err
				}
				value, err := etcdCnfx.GetObject(t.Key)
				if err != nil {
					return err
				}
				if t.CB != nil {
					err := t.Init(value)
					if err != nil {
						return err
					}
					if t.Watch {
						etcdCnfxReader.Watch(context.TODO(), t.Key, func(etcdValue *config_etcd.EtcdValue, err error) {
							if err != nil {
								return
							}
							t.Init(value)
						})
					}
				}
			}
		}
	}
	return etcdCnfxReader.Close()
}
