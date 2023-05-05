package http

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/url"
)

type IHttpReuqest interface {
	Url() *url.URL
	ContentType() string
	Headers() IWebHeaderCollection
	Method() Method
	Reader() io.ReadCloser
	Hash() (string, error)
}

type httpRequest struct {
	url     *url.URL
	headers IWebHeaderCollection
	method  Method
	reader  io.ReadCloser
}

func NewHttpRequest(url *url.URL, method Method, contentType string, headers IWebHeaderCollection, reader io.Reader) IHttpReuqest {
	httpRequest := httpRequest{
		url:     url,
		headers: headers,
		method:  method,
		reader:  io.NopCloser(reader),
	}
	return &httpRequest
}

func (httpRequest httpRequest) Url() *url.URL {
	return httpRequest.url
}

func (httpRequest httpRequest) ContentType() string {
	value, _ := httpRequest.headers.Get("content-type")
	return value
}

func (httpRequest httpRequest) Headers() IWebHeaderCollection {
	return httpRequest.headers
}

func (httpRequest httpRequest) Method() Method {
	return httpRequest.method
}

func (httpRequest httpRequest) Reader() io.ReadCloser {
	if httpRequest.method == GET {
		return io.NopCloser(bytes.NewBuffer([]byte{}))
	}
	return httpRequest.reader
}

func (httpRequest *httpRequest) Hash() (string, error) {
	buffer := bytes.NewBufferString("")
	buffer.WriteString(httpRequest.url.String())
	buffer.WriteString("#")
	var b bytes.Buffer
	_, err := io.Copy(&b, httpRequest.Reader())
	if err != nil {
		return "", err
	}
	httpRequest.reader = io.NopCloser(bytes.NewReader(b.Bytes()))
	buffer.WriteString(b.String())
	sha256 := sha256.New()
	_, err = sha256.Write(buffer.Bytes())
	if err != nil {
		return "", err
	}
	hash := sha256.Sum(nil)
	return hex.EncodeToString(hash), nil
}
