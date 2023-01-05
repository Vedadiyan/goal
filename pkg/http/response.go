package http

import (
	"io"
	"net/http"
	"strings"
)

type IHttpResponse interface {
	Status() int
	Headers() http.Header
	ContentType() string
	Reader() io.ReadCloser
}

type httpResponse struct {
	response http.Response
}

func (httpResponse httpResponse) Status() int {
	return httpResponse.response.StatusCode
}

func (httpResponse httpResponse) Reader() io.ReadCloser {
	return httpResponse.response.Body
}

func (httpResponse httpResponse) Headers() http.Header {
	return httpResponse.response.Header
}

func (httpResponse httpResponse) ContentType() string {
	return strings.ToLower(strings.Split(httpResponse.response.Header.Get("Content-Type"), ";")[0])
}
