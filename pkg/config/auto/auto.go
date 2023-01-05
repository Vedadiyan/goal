package config_auto

import (
	"github.com/vedadiyan/goal/pkg/config"
	config_etcd "github.com/vedadiyan/goal/pkg/config/etcd"
)

type IInitializer interface {
	Inititialize(value any) error
}

type String struct {
	Key  string
	Init func(value string)
}

type KeyValue map[string]any

type Object struct {
	Key  string
	Init func(value KeyValue)
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

func (s String) Inititialize(value any) error {
	if str, ok := value.(string); ok {
		s.Init(str)
		return nil
	}
	return config.INVALID_OBJECT
}
func (o Object) Inititialize(value any) error {
	if keyValue, ok := value.(map[string]any); ok {
		o.Init(keyValue)
		return nil
	}
	return config.INVALID_OBJECT
}

func Bootstrap(url string, key string, initializers ...IInitializer) error {
	etcdCnfxReader := config_etcd.New(url, key)
	etcdCnfx, err := etcdCnfxReader.ReadConfig()
	if err != nil {
		return err
	}
	for _, initializer := range initializers {
		switch t := initializer.(type) {
		case String:
			{
				value, err := etcdCnfx.GetString(t.Key)
				if err != nil {
					return err
				}
				if t.Init != nil {
					err := t.Inititialize(value)
					if err != nil {
						return err
					}
				}
			}
		case Object:
			{
				value, err := etcdCnfx.GetObject(t.Key)
				if err != nil {
					return err
				}
				if t.Init != nil {
					err := t.Inititialize(value)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
