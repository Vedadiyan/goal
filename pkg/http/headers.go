package http

import (
	"fmt"
	"net/http"
	"strings"
)

type IWebHeaderCollection interface {
	Add(key string, value string)
	Get(key string) (string, error)
	GetOrDefault(key string, defaultValue string) string
	Range() map[string]string
	Copy(header *http.Header)
}

type webHeaderCollection struct {
	headers map[string]string
}

func NewWebHeaderCollection() IWebHeaderCollection {
	webHeaderCollection := &webHeaderCollection{
		headers: make(map[string]string),
	}
	return webHeaderCollection
}

func (webHeaderCollection *webHeaderCollection) Add(key string, value string) {
	webHeaderCollection.headers[strings.ToLower(key)] = value
}

func (webHeaderCollection webHeaderCollection) Get(key string) (string, error) {
	if value, ok := webHeaderCollection.headers[strings.ToLower(key)]; ok {
		return value, nil
	}
	return "", fmt.Errorf("key not found")
}

func (webHeaderCollection webHeaderCollection) GetOrDefault(key string, defaultValue string) string {
	if value, ok := webHeaderCollection.headers[strings.ToLower(key)]; ok {
		return value
	}
	return defaultValue
}

func (webHeaderCollection webHeaderCollection) Range() map[string]string {
	return webHeaderCollection.headers
}

func (webHeaderCollection webHeaderCollection) Copy(header *http.Header) {
	for key, value := range webHeaderCollection.Range() {
		header.Add(key, value)
	}
}
