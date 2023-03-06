package config_etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vedadiyan/goal/pkg/config"
	"github.com/vedadiyan/goal/pkg/util"
	etcdclient "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	client *etcdclient.Client
}

type EtcdValue struct {
	data map[string]any
}

func (e EtcdValue) GetString(key string) (string, error) {
	if key == "" {
		key = "$"
	}
	ref, ok := e.data[key]
	if !ok {
		return "", config.KEY_NOT_FOUND
	}
	value, ok := ref.(string)
	if !ok {
		return "", config.INVALID_CAST
	}
	return value, nil
}
func (e EtcdValue) GetNumber(key string) (float64, error) {
	if key == "" {
		key = "$"
	}
	ref, ok := e.data[key]
	if !ok {
		return 0, config.KEY_NOT_FOUND
	}
	value, ok := ref.(float64)
	if !ok {
		return 0, config.INVALID_CAST
	}
	return value, nil
}
func (e EtcdValue) GetBoolean(key string) (bool, error) {
	if key == "" {
		key = "$"
	}
	ref, ok := e.data[key]
	if !ok {
		return false, config.KEY_NOT_FOUND
	}
	value, ok := ref.(bool)
	if !ok {
		return false, config.INVALID_CAST
	}
	return value, nil
}

func (e EtcdValue) GetObject(key string) (map[string]any, error) {
	_key := key + "."
	m := make(map[string]any)
	for k, value := range e.data {
		if strings.HasPrefix(k, _key) {
			m[k[len(_key):]] = value
		}
	}
	return m, nil
}

func NewClient(endpoints []string) (*EtcdClient, error) {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, err
	}
	etcdClient := EtcdClient{
		client: client,
	}
	return &etcdClient, nil
}

func (e EtcdClient) Close() error {
	return e.client.Close()
}

func (e EtcdClient) ReadKey(ctx context.Context, key string) (*EtcdValue, error) {
	value, err := e.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if value.Count == 0 {
		return nil, fmt.Errorf("key not found")
	}
	if value.Count > 1 {
		return nil, fmt.Errorf("multiple keys found")
	}
	strValue := string(value.Kvs[0].Value)
	return NewEtcdValue(key, strValue)
}

func NewEtcdValue(key string, value string) (*EtcdValue, error) {
	etcdValue := EtcdValue{}
	if util.IsJSON(value) {
		mapper := make(map[string]any)
		err := json.Unmarshal([]byte(value), &mapper)
		if err != nil {
			return nil, err
		}
		mapper = flattenValue(&key, mapper)
		etcdValue.data = mapper
		return &etcdValue, nil
	}
	etcdValue.data = map[string]any{
		"$": value,
	}
	return &etcdValue, nil
}

func (e EtcdClient) Watch(ctx context.Context, key string, fn func(etcdValue *EtcdValue, err error)) {
	watcher := e.client.Watch(ctx, key)
	go func() {
		for value := range watcher {
			if value.Events == nil || value.Events[0] == nil {
				fn(nil, config.INVALID_RESPONSE)
				continue
			}
			fn(NewEtcdValue(string(value.Events[0].Kv.Key), string(value.Events[0].Kv.Value)))
		}
	}()
}

func flattenValue(key *string, data map[string]any) map[string]any {
	output := make(map[string]any)
	for k, v := range data {
		innerValue, ok := v.(map[string]any)
		if ok {
			innerMap := flattenValue(key, innerValue)
			for innerKey, innerValue := range innerMap {
				tmp := k + "." + innerKey
				output[tmp] = innerValue
			}
		} else {
			output[k] = v
		}
	}
	return output
}
