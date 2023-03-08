package http

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/url"
	"strings"
)

type JSON string
type XML string
type URLEncoded = url.Values

func Send(url IUrl, defaultHeaders IWebHeaderCollection, method Method, request any) (res IHttpResponse, err error) {
	return SendWithContext(context.TODO(), url, defaultHeaders, method, request)
}
func SendWithContext(ctx context.Context, url IUrl, defaultHeaders IWebHeaderCollection, method Method, request any) (IHttpResponse, error) {
	rqUrl, err := url.Url()
	if err != nil {
		return nil, err
	}
	headers := defaultHeaders
	if headers == nil {
		headers = NewWebHeaderCollection()
	}
	rqType, readCloser := GetRequest(request)
	rq := httpRequest{
		url:         rqUrl,
		contentType: rqType,
		headers:     headers,
		method:      method,
		reader:      readCloser,
	}
	response, err := GetHttpClient().Send(ctx, &rq)
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

func Nil() io.ReadCloser {
	return io.NopCloser(bytes.NewReader([]byte{}))
}

func GetRequest(request any) (string, io.ReadCloser) {
	switch t := request.(type) {
	case string:
		{
			return "text/plain", io.NopCloser(strings.NewReader(t))
		}
	case JSON:
		{
			return "application/json", io.NopCloser(strings.NewReader(string(t)))
		}
	case XML:
		{
			return "text/xml", io.NopCloser(strings.NewReader(string(t)))
		}
	case URLEncoded:
		{
			return "application/x-www-form-urlencoded", io.NopCloser(strings.NewReader(t.Encode()))
		}
	case []byte:
		{
			return "application/octet-stream", io.NopCloser(bytes.NewReader(t))
		}
	}
	return "", nil
}
