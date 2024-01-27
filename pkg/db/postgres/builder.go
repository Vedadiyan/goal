package postgres

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"text/template"

	"github.com/vedadiyan/goal/pkg/db/postgres/sanitize"
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
		template *template.Template
		err      error
	)
	hash, err := templateHash(str)
	rwmux.RLock()
	template, ok := templates[hash]
	rwmux.RUnlock()
	if !ok {
		template, err = template.New(hash).Parse(str)
		if err != nil {
			return "", err
		}
		template.Funcs(map[string]any{
			"Sanitize": func(value any) string {
				val, _err := sanitize.SanitizeSQL("$1", standardize(value))
				err = _err
				return val
			},
		})
		rwmux.Lock()
		templates[hash] = template
		rwmux.Unlock()
	}
	if err != nil {
		return "", err
	}
	buffer := bytes.NewBufferString("")
	err = template.Execute(buffer, args)
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
			return int32(value)
		}
	case uint:
		{
			return uint32(value)
		}
	default:
		{
			return value
		}
	}
}
