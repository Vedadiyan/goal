package config_etcd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vedadiyan/goal/pkg/config"
)

type EtcdConfig struct {
	url  string
	data map[string]any
	key  string
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

func (e *EtcdConfig) ReadConfig() (config.ObjectReader, error) {
	base64Key := base64.StdEncoding.EncodeToString([]byte(e.key))
	base64KeyBytes := bytes.NewBufferString(fmt.Sprintf(`{ "key": "%s" }`, base64Key))
	res, err := http.Post(e.url+"v3/kv/range", "application/json", base64KeyBytes)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 200 {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		data := make(map[string]any)
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, err
		}
		if data["kvs"] == nil {
			return nil, config.KEY_NOT_FOUND
		}
		kvsArray, ok := data["kvs"].([]any)
		if !ok {
			return nil, config.INVALID_CAST
		}
		kvsMap, ok := kvsArray[0].(map[string]any)
		if !ok {
			return nil, config.INVALID_CAST
		}
		if kvsMap["value"] == nil {
			return nil, config.NO_VALUE
		}
		kvsValueBase64, ok := kvsMap["value"].(string)
		if !ok {
			return nil, config.INVALID_CAST
		}
		kvsValue, err := base64.StdEncoding.DecodeString(kvsValueBase64)
		if err != nil {
			return nil, err
		}
		nonFlatData := make(map[string]any)
		if err := json.Unmarshal(kvsValue, &nonFlatData); err != nil {
			return nil, err
		}
		e.data = flattenValue(nil, nonFlatData)
		return e, nil
	}
	return nil, config.INVALID_RESPONSE
}

func (e EtcdConfig) GetString(key string) (string, error) {
	if e.data[key] == nil {
		return "", config.KEY_NOT_FOUND
	}
	value, ok := e.data[key].(string)
	if !ok {
		return "", config.INVALID_CAST
	}
	return value, nil
}
func (e EtcdConfig) GetInt(key string) (int, error) {
	if e.data[key] == nil {
		return 0, config.KEY_NOT_FOUND
	}
	value, ok := e.data[key].(int)
	if !ok {
		return 0, config.INVALID_CAST
	}
	return value, nil
}
func (e EtcdConfig) GetInt64(key string) (int64, error) {
	if e.data[key] == nil {
		return 0, config.KEY_NOT_FOUND
	}
	value, ok := e.data[key].(int64)
	if !ok {
		return 0, config.INVALID_CAST
	}
	return value, nil
}
func (e EtcdConfig) GetBoolean(key string) (bool, error) {
	if e.data[key] == nil {
		return false, config.KEY_NOT_FOUND
	}
	value, ok := e.data[key].(bool)
	if !ok {
		return false, config.INVALID_CAST
	}
	return value, nil
}

func (e EtcdConfig) GetObject(key string) (map[string]any, error) {
	_key := key + "."
	m := make(map[string]any)
	for k, value := range e.data {
		if strings.HasPrefix(k, _key) {
			m[k[len(_key):]] = value
		}
	}
	return m, nil
}

func New(url string, key string) config.ConfigReader {
	confix := EtcdConfig{
		url: url,
		key: key,
	}
	return &confix
}
