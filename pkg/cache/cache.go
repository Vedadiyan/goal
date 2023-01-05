package cache

import (
	"sync"
	"time"
)

var store sync.Map

func init() {
	store = sync.Map{}
}

func Add(key string, value any) error {
	if _, ok := store.LoadOrStore(key, value); ok {
		return DUPLICATE_KEY
	}
	return nil
}

func AddWithTTL(key string, value any, ttl time.Duration) error {
	err := Add(key, value)
	if err != nil {
		return err
	}
	time.AfterFunc(ttl, func() {
		Delete(key)
	})
	return nil
}

func Delete(key string) error {
	if _, ok := store.LoadAndDelete(key); !ok {
		return KEY_NOT_FOUND
	}
	return nil
}

func Set(key string, value any) error {
	if _, ok := store.Load(key); !ok {
		return KEY_NOT_FOUND
	}
	store.Store(key, value)
	return nil
}

func Get[T any](key string) (T, error) {
	value, ok := store.Load(key)
	if !ok {
		return *new(T), KEY_NOT_FOUND
	}
	v, ok := value.(T)
	if !ok {
		return *new(T), INVALID_CAST
	}
	return v, nil
}

func GetOrSet[T any](key string, value any) (T, error) {
	store.LoadOrStore(key, value)
	v, ok := value.(T)
	if !ok {
		return *new(T), INVALID_CAST
	}
	return v, nil
}
