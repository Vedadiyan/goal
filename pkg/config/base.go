package config

type ObjectReader interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetInt64(key string) (int64, error)
	GetBoolean(key string) (bool, error)
	GetObject(key string) (map[string]any, error)
}

type ConfigReader interface {
	ReadConfig() (ObjectReader, error)
}
