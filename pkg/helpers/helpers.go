package helpers

func ToPtr[T any](value T) *T {
	return &value
}

func UnsafeValue[T any](ptr *T) T {
	return *ptr
}
