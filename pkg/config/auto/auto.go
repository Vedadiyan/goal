package config_auto

import (
	"context"
	"encoding/json"
	"os"
	"strings"
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
type ConfigMapBootstrapper struct{}
type ConfigBootstrapper interface {
	Bootstrap() error
}
type ETCDBootstrapper struct {
	url string
}

var (
	_initializers []any
	rwMut         sync.RWMutex
)

func init() {
	_initializers = make([]any, 0)
}

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
	rwMut.Lock()
	defer rwMut.Unlock()
	_initializers = append(_initializers, initializer)
}
func (configMapBootstrapper *ConfigMapBootstrapper) Bootstrap() error {
	env := make(map[string]string)
	for _, item := range os.Environ() {
		keyValue := strings.SplitN(item, "=", 2)
		env[keyValue[0]] = keyValue[1]

	}

	for _, value := range _initializers {
		initializer := value.(Initializer)
		switch t := initializer.(type) {
		case String:
			{
				value, ok := env[t.Key]
				if !ok {
					return config.KEY_NOT_FOUND
				}
				if t.CB != nil {
					err := t.Init(value)
					if err != nil {
						return err
					}
				}
			}
		case Object:
			{
				value, ok := env[t.Key]
				if !ok {
					return config.KEY_NOT_FOUND
				}
				objectValue := make(map[string]any)
				err := json.Unmarshal([]byte(value), &objectValue)
				if err != nil {
					return err
				}
				if t.CB != nil {
					err := t.Init(value)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
func (etcdBootstrapper *ETCDBootstrapper) Bootstrap() error {
	etcdCnfxReader, err := config_etcd.NewClient([]string{etcdBootstrapper.url})
	if err != nil {
		return err
	}
	for _, value := range _initializers {
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
							_ = t.Init(value)
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
							_ = t.Init(value)
						})
					}
				}
			}
		}
	}
	return etcdCnfxReader.Close()
}
func ForETCD(url string) *ETCDBootstrapper {
	etcdBootstrapper := &ETCDBootstrapper{
		url: url,
	}
	return etcdBootstrapper
}
func ForConfigMap() *ConfigMapBootstrapper {
	configMapBootstrapper := &ConfigMapBootstrapper{}
	return configMapBootstrapper
}
func Bootstrap(configBootstrapper ConfigBootstrapper) error {
	return configBootstrapper.Bootstrap()
}
func New[T string | KeyValue](key string, watch bool, cb func(value T)) Initializer {
	var value T
	switch any(value).(type) {
	case string:
		{
			initializer := String{
				Key:   key,
				Watch: watch,
				CB: func(value string) {
					cb(any(value).(T))
				},
			}
			return initializer
		}
	case KeyValue:
		{
			Initializer := Object{
				Key:   key,
				Watch: watch,
				CB: func(value KeyValue) {
					cb(any(value).(T))
				},
			}
			return Initializer
		}
	}
	panic("unexpected case")
}
