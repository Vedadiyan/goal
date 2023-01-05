package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"sync"
	"time"

	"github.com/vedadiyan/goal/pkg/cache"
)

type Option any

type CacheOption int

const (
	FETCH_OPTION_CACHED CacheOption = 1
)

var (
	ConfigureMarshallerMethods func(t any) (MarshallerType, string, any)
	IsDebug                    bool = false
	_once                      sync.Once
	_httpClient                IHttpClient
	_cachFn                    CacheFn
	_defaultTTL                time.Duration
)

type CacheFn func(url *url.URL, key string, value func() (IHttpResponse, error), ttl time.Duration) (IHttpResponse, error)

type IHttpClient interface {
	Send(ctx context.Context, httpRequest IHttpReuqest, options ...Option) (IHttpResponse, error)
}

type httpClient struct {
	httpClient http.Client
}

func init() {
	_defaultTTL = time.Second * 30
	_cachFn = func(url *url.URL, key string, value func() (IHttpResponse, error), ttl time.Duration) (IHttpResponse, error) {
		cahcedValue, err := cache.Get[IHttpResponse](key)
		if errors.Is(err, cache.KEY_NOT_FOUND) {
			res, err := value()
			if err != nil {
				return nil, err
			}
			cache.AddWithTTL(key, res, ttl)
			return res, nil
		}
		return cahcedValue, nil
	}
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
	if hasCacheOption(FETCH_OPTION_CACHED, options) {
		return readOrSend(&httpClient, ctx, httpRequest, getTTL(options))
	}
	return send(&httpClient, ctx, httpRequest)
}

func readOrSend(httpClient *httpClient, ctx context.Context, httpRequest IHttpReuqest, ttl time.Duration) (IHttpResponse, error) {
	hash, err := httpRequest.Hash()
	if err != nil {
		return nil, err
	}
	return _cachFn(httpRequest.Url(), hash, func() (IHttpResponse, error) {
		return send(httpClient, ctx, httpRequest)
	}, ttl)
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

func hasCacheOption(option CacheOption, options []Option) bool {
	for _, item := range options {
		if value, ok := item.(CacheOption); ok {
			if value&option == option {
				return true
			}
		}
	}
	return false
}

func getTTL(options []Option) time.Duration {
	for _, item := range options {
		if value, ok := item.(time.Duration); ok {
			return value
		}
	}
	return _defaultTTL
}

func RegisterCache(cacheFn CacheFn) {
	_cachFn = cacheFn
}

func ConfigureDefaultCacheTTL(defaultTTL time.Duration) {
	_defaultTTL = defaultTTL
}
