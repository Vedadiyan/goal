package util

import "regexp"

const (
	_JSON_PATTERN = `^(\{(\s*"[^"]+"\s*:\s*(?:"[^"]*"|\d+|true|false|null)(?:\s*,\s*"[^"]+"\s*:\s*(?:"[^"]*"|\d+|true|false|null))*\s*|)\}|\[(\s*(?:"[^"]*"|\d+|true|false|null)(?:\s*,\s*(?:"[^"]*"|\d+|true|false|null))*\s*|)\])$`
)

var (
	_jsonRegexp regexp.Regexp
)

func init() {
	jsonRegexp, err := regexp.Compile(_JSON_PATTERN)
	if err != nil {
		panic(err)
	}
	_jsonRegexp = *jsonRegexp
}

func IsJSON(str string) bool {
	return _jsonRegexp.Match([]byte(str))
}
