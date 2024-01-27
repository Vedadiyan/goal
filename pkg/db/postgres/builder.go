package postgres

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/vedadiyan/goal/pkg/db/postgres/sanitize"
)

type (
	TemplateData struct {
		Root map[string]any
	}
)

var (
	templates map[string]*template.Template
	rwmux     sync.RWMutex
)

func init() {
	templates = make(map[string]*template.Template)
}

func Build(str string, args map[string]any) (string, error) {
	var (
		_template *template.Template
		err       error
	)
	hash, err := templateHash(str)
	rwmux.RLock()
	_template, ok := templates[hash]
	rwmux.RUnlock()
	if !ok {
		_template = template.New(hash)
		_template.Funcs(map[string]any{
			"Sanitize": func(value any) string {
				val, _err := sanitize.SanitizeSQL("$1", standardize(value))
				err = _err
				return val
			},
			"DateRange": func(from string, to string) string {
				_from, _err := sanitize.SanitizeSQL("$1", standardize(from))
				if err != nil {
					err = _err
					return ""
				}
				_to, _err := sanitize.SanitizeSQL("$1", standardize(to))
				if err != nil {
					err = _err
					return ""
				}
				return fmt.Sprintf("'[%s, %s]'::daterange", strings.ReplaceAll(_from, "'", ""), strings.ReplaceAll(_to, "'", ""))
			},
		})
		_template, err := _template.Parse(str)
		if err != nil {
			return "", err
		}
		rwmux.Lock()
		templates[hash] = _template
		rwmux.Unlock()
	}
	if err != nil {
		return "", err
	}
	data := TemplateData{
		Root: args,
	}
	buffer := bytes.NewBufferString("")
	err = _template.Execute(buffer, data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func templateHash(str string) (string, error) {
	sha256 := sha256.New()
	_, err := sha256.Write([]byte(str))
	if err != nil {
		return "", err
	}
	hash := sha256.Sum(nil)
	return hex.EncodeToString(hash), nil
}

func standardize(value any) any {
	switch value := value.(type) {
	case int:
		{
			return int64(value)
		}
	case int32:
		{
			return int64(value)
		}
	case int16:
		{
			return int64(value)
		}
	case int8:
		{
			return int64(value)
		}
	case uint:
		{
			return uint64(value)
		}
	case uint32:
		{
			return uint64(value)
		}
	case uint16:
		{
			return uint64(value)
		}
	case uint8:
		{
			return uint64(value)
		}
	default:
		{
			return value
		}
	}
}
