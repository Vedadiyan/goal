package config

type ConfigError string

const (
	KEY_NOT_FOUND    ConfigError = ConfigError("key not found")
	INVALID_CAST     ConfigError = ConfigError("invalid cast")
	NO_VALUE         ConfigError = ConfigError("no value available")
	INVALID_OBJECT   ConfigError = ConfigError("invalid object")
	INVALID_RESPONSE ConfigError = ConfigError("etcd returned an invalid response")
)

func (c ConfigError) Error() string {
	return string(c)
}
