package http

type Method string
type MarshallerType int

const (
	NONE MarshallerType = iota
	JSON
	PROTOBUF
	URL_ENCODED_FORM
	BINARY
)

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATH   Method = "PATCH"
	DELETE Method = "DELETE"
)
