package http

import (
	"bytes"
	"context"
	"io"
	"log"
)

type NEVER *byte

func Send(url IUrl, defaultHeaders IWebHeaderCollection, method Method, request io.ReadCloser) (res IHttpResponse, err error) {
	rqUrl, err := url.Url()
	if err != nil {
		return nil, err
	}
	headers := defaultHeaders
	if headers == nil {
		headers = NewWebHeaderCollection()
	}
	rq := httpRequest{
		url:         rqUrl,
		contentType: headers.GetOrDefault("Content-Type", "text/plain"),
		headers:     headers,
		method:      method,
		reader:      request,
	}
	response, err := GetHttpClient().Send(context.TODO(), &rq)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func Must[TType any](fn func() (TType, error)) TType {
	res, err := fn()
	if err != nil {
		log.Fatalln(err)
	}
	return res
}

func Never() NEVER {
	return nil
}

func Nil() io.ReadCloser {
	return io.NopCloser(bytes.NewReader([]byte{}))
}
