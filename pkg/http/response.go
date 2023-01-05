package http

import "io"

type IHttpResponse interface {
	Status() int
	Reader() io.ReadCloser
}

type httpResponse struct {
	status int
	reader io.ReadCloser
}

func (httpResponse httpResponse) Status() int {
	return httpResponse.status
}

func (httpResponse httpResponse) Reader() io.ReadCloser {
	return httpResponse.reader
}
