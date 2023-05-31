package protoval

import (
	"fmt"
	"strconv"
)

var (
	_rules map[string]func(name string, value any, rule string) error
)

func init() {
	_rules = make(map[string]func(name string, value any, rule string) error)
	_rules["required"] = Required
	_rules["max_len"] = MaxLen
	_rules["min_len"] = MinLen
}

func Register(name string, fn func(name string, value any, rule string) error) {
	_rules[name] = fn
}

func Required(name string, value any, rule string) error {
	if value == nil {
		return Error(name, "is required")
	}
	return nil
}

func MinLen(name string, value any, rule string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string but recieved %T", value)
	}
	i, err := strconv.ParseInt(rule, 10, 32)
	if err != nil {
		return err
	}
	if len(str) < int(i) {
		return Error(name, fmt.Sprintf("must be larger than or equal in size to %d", i))
	}
	return nil
}

func MaxLen(name string, value any, rule string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string but recieved %T", value)
	}
	i, err := strconv.ParseInt(rule, 10, 32)
	if err != nil {
		return err
	}
	if len(str) > int(i) {
		return Error(name, fmt.Sprintf("must be smaller than or equal in size to %d", i))
	}
	return nil
}
