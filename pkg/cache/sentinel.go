package cache

type CacheError string

const (
	KEY_NOT_FOUND CacheError = CacheError("key not found")
	DUPLICATE_KEY CacheError = CacheError("an object with the same key already has already been added")
	INVALID_CAST  CacheError = CacheError("invalid cast")
)

func (c CacheError) Error() string {
	return string(c)
}
