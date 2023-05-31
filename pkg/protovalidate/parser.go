package protoval

import (
	"bytes"
	"fmt"
	"strings"
)

func TrimLeftRight(str string) string {
	return strings.TrimRight(strings.TrimLeft(str, " "), " ")
}

func GetRule(str string, rules map[string]string) error {
	if len(str) == 0 {
		return nil
	}
	split := strings.Split(str, ":")
	if len(split) == 1 {
		rules[str] = "true"
		return nil
	}
	if len(split) == 2 {
		rules[strings.TrimRight(split[0], " ")] = strings.TrimLeft(split[1], " ")
		return nil
	}
	return fmt.Errorf("invalid rule")
}

func ExprParser(expr string) (map[string]string, error) {
	var buffer bytes.Buffer
	rules := make(map[string]string)

	for i := 0; i < len(expr); i++ {
		c := expr[i]
		switch c {
		case '/':
			{
				if buffer.Len() > 0 {
					return nil, fmt.Errorf("invalid pattern")
				}
				hold := false
			LOOP:
				for x := i + 1; x < len(expr); x++ {
					c := expr[x]
					if hold {
						hold = false
						goto WRITE
					}
					switch c {
					case '\\':
						{
							hold = true
						}
					case '/':
						{
							i = x
							break LOOP
						}
					}
				WRITE:
					buffer.WriteByte(c)
				}
				str := buffer.String()
				buffer.Reset()
				rules["regex"] = str
			}
		case '|':
			{
				str := TrimLeftRight(buffer.String())
				buffer.Reset()
				err := GetRule(str, rules)
				if err != nil {
					return nil, err
				}
			}
		default:
			{
				buffer.WriteByte(c)
			}
		}
	}
	str := TrimLeftRight(buffer.String())
	buffer.Reset()
	err := GetRule(str, rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}
