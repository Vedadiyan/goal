package http

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"strings"
)

type JSON string
type XML string
type URLEncoded = url.Values

func Send[T string | JSON | XML | URLEncoded | []byte](url IUrl, defaultHeaders IWebHeaderCollection, method Method, request T) (res IHttpResponse, err error) {
	return SendWithContext(context.TODO(), url, defaultHeaders, method, request)
}
func SendWithContext[T string | JSON | XML | URLEncoded | []byte](ctx context.Context, url IUrl, defaultHeaders IWebHeaderCollection, method Method, request T) (IHttpResponse, error) {
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

func Nil() []byte {
	return nil
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
			return "application/xml", io.NopCloser(strings.NewReader(string(t)))
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
