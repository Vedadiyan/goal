package http

import (
	"context"
	"log"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"
)

type Option any

type CacheOption int

var (
	ConfigureMarshallerMethods func(t any) (MarshallerType, string, any)
	IsDebug                    bool = false
	_once                      sync.Once
	_httpClient                IHttpClient
)

type IHttpClient interface {
	Send(ctx context.Context, httpRequest IHttpReuqest, options ...Option) (IHttpResponse, error)
}

type httpClient struct {
	httpClient http.Client
}

func ConfigureHttpClient(transport *http.Transport, timeout time.Duration) {
	_once.Do(func() {
		_httpClient = &httpClient{
			httpClient: http.Client{
				Transport: transport,
				Timeout:   timeout,
			},
		}
	})
}

func GetHttpClient() IHttpClient {
	_once.Do(func() {
		_httpClient = &httpClient{
			httpClient: http.Client{
				Transport: &http.Transport{
					MaxIdleConns:        100,
					MaxIdleConnsPerHost: 20,
				},
			},
		}
	})
	return _httpClient
}

func (httpClient httpClient) Send(ctx context.Context, httpRequest IHttpReuqest, options ...Option) (IHttpResponse, error) {
	return send(&httpClient, ctx, httpRequest)
}

func send(httpClient *httpClient, ctx context.Context, httpRequest IHttpReuqest) (IHttpResponse, error) {
	url := httpRequest.Url()
	var request *http.Request
	var err error
	if !IsDebug {
		request, err = http.NewRequestWithContext(ctx, string(httpRequest.Method()), url.String(), httpRequest.Reader())
	} else {
		request, err = http.NewRequestWithContext(httptrace.WithClientTrace(ctx, debugConnectionReuse()), string(httpRequest.Method()), url.String(), httpRequest.Reader())
	}
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	if httpRequest.ContentType() != "" {
		request.Header.Add("Content-Type", httpRequest.ContentType())
	}
	if httpRequest.Headers() != nil {
		httpRequest.Headers().Copy(&request.Header)
	}
	response, err := httpClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return &httpResponse{response: *response}, nil
}

func debugConnectionReuse() *httptrace.ClientTrace {
	clientTrace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			log.Printf("%s wasIdle: %t reused: %t idleTime: %d", info.Conn.RemoteAddr().String(), info.WasIdle, info.Reused, info.IdleTime)
		},
	}
	return clientTrace
}
