package cache

import (
	"runtime"
	"sync"
	"time"
)

type Events int

const (
	ADD Events = iota
	SET
	DELETE
)

var (
	_store    sync.Map
	_watchers sync.Map
)

func init() {
	_store = sync.Map{}
}

func Add(key string, value any) error {
	if _, ok := _store.LoadOrStore(key, value); ok {
		return DUPLICATE_KEY
	}
	raise(ADD, key, value)
	return nil
}

func AddWithTTL(key string, value any, ttl time.Duration) error {
	err := Add(key, value)
	if err != nil {
		return err
	}
	time.AfterFunc(ttl, func() {
		_ = Delete(key)
	})
	raise(ADD, key, value)
	return nil
}

func Delete(key string) error {
	value, ok := _store.LoadAndDelete(key)
	if !ok {
		return KEY_NOT_FOUND
	}
	raise(DELETE, key, value)
	return nil
}

func Set(key string, value any) error {
	_store.Store(key, value)
	raise(SET, key, value)
	return nil
}

func SetWithTTL(key string, value any, ttl time.Duration) error {
	_store.Store(key, value)
	var timer *time.Timer
	timer = time.AfterFunc(ttl, func() {
		_ = Delete(key)
		runtime.KeepAlive(timer)
	})
	raise(SET, key, value)
	return nil
}

func Get[T any](key string) (T, error) {
	value, ok := _store.Load(key)
	if !ok {
		return *new(T), KEY_NOT_FOUND
	}
	v, ok := value.(T)
	if !ok {
		return *new(T), INVALID_CAST
	}
	return v, nil
}

func Watch(key string, cb func(Events, any)) {
	value, ok := _watchers.Load(key)
	if !ok {
		value = make([]func(Events, any), 0)
	}
	value = append(value.([]func(Events, any)), cb)
	_watchers.Store(key, value)
}

func raise(event Events, key string, value any) {
	values, ok := _watchers.Load(key)
	if !ok {
		return
	}
	for _, value := range values.([]func(Events, any)) {
		value(event, value)
	}
}
