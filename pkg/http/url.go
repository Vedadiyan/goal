package http

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
)

type Query map[string][]string
type RouteValues map[string]string

type IUrl interface {
	Url() (*url.URL, error)
}

type UrlTemplate struct {
	template string
}

func NewUrlTemplate(template string) *UrlTemplate {
	urlTemplate := UrlTemplate{
		template: template,
	}
	return &urlTemplate
}

func (urlTemplate UrlTemplate) Url(routeValues RouteValues, query Query) (*url.URL, error) {
	copy := urlTemplate.template
	buffer := bytes.NewBufferString("?")
	for key, value := range routeValues {
		copy = strings.ReplaceAll(copy, key, value)
	}
	for key, values := range query {
		for _, value := range values {
			buffer.WriteString(fmt.Sprintf("%s=%s&", key, value))
		}
	}
	bytes := buffer.Bytes()
	finalUrl := fmt.Sprintf("%s%s", copy, string(bytes[:len(bytes)-1]))
	return url.Parse(finalUrl)
}
