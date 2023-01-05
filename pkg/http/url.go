package http

import (
	"bytes"
	"fmt"
	"net/url"
)

type Query map[string]string

type IUrl interface {
	Url() (*url.URL, error)
}

type urlBuilder struct {
	relativeUrl string
	query       Query
}

func (urlBuilder urlBuilder) Url() (*url.URL, error) {
	var _bytes []byte
	buffer := bytes.NewBufferString(urlBuilder.relativeUrl)
	if urlBuilder.query != nil {
		buffer.WriteString("?")
		for key, value := range urlBuilder.query {
			buffer.WriteString(fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
			buffer.WriteString("&")
		}
		_bytes = buffer.Bytes()
		_bytes = _bytes[0 : len(_bytes)-1]
	} else {
		_bytes = buffer.Bytes()
	}
	return url.Parse(string(_bytes))
}

func Url(url string, query Query) *urlBuilder {
	return &urlBuilder{url, query}
}

func BaseUrl(baseUrl string) *urlBuilder {
	return &urlBuilder{baseUrl, nil}
}

func ExtendUrl(baseUrl *url.URL, relativeUrl string, query Query) *urlBuilder {
	return &urlBuilder{fmt.Sprintf("%s://%s/%s", baseUrl.Scheme, baseUrl.Host, relativeUrl), query}
}
