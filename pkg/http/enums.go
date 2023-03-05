package http

type Method string
type MarshallerType int

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATH   Method = "PATCH"
	DELETE Method = "DELETE"
)
