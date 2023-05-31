package protoval

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	_rules map[string]func(name string, value any, rule string) error
)

func init() {
	_rules = make(map[string]func(name string, value any, rule string) error)
	_rules["required"] = Required
	_rules["max_len"] = MaxLen
	_rules["min_len"] = MinLen
	_rules["latitude"] = Latitude
	_rules["longitude"] = Longitude
	_rules["email"] = Email
	_rules["future_date"] = FutureDate
	_rules["mix"] = Min
	_rules["max"] = Max
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

func Latitude(name string, value any, rule string) error {
	val, ok := value.(float64)
	if !ok {
		return fmt.Errorf("expected float64 but recieved %T", value)
	}
	if val < -90 || val > 90 {
		return Error(name, "invalid latitude")
	}
	return nil
}

func Longitude(name string, value any, rule string) error {
	val, ok := value.(float64)
	if !ok {
		return fmt.Errorf("expected float64 but recieved %T", value)
	}
	if val < -180 || val > 180 {
		return Error(name, "invalid longitude")
	}
	return nil
}

func Email(name string, value any, rule string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string but recieved %T", value)
	}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	ok, _ = regexp.MatchString(pattern, str)
	if !ok {
		return Error(name, "invalid email")
	}
	return nil
}

func FutureDate(name string, value any, rule string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string but recieved %T", value)
	}
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		return Error(name, "inavlid date")
	}
	if !date.After(time.Now()) {
		return Error(name, "must be equal to or after the current date")
	}
	return nil
}

func Min(name string, value any, rule string) error {
	val, err := getNumber(value)
	if err != nil {
		return err
	}
	i, err := strconv.ParseFloat(rule, 64)
	if err != nil {
		return err
	}
	if val > i {
		return Error(name, fmt.Sprintf("must be less than or equal to %f", i))
	}
	return nil
}

func Max(name string, value any, rule string) error {
	val, err := getNumber(value)
	if err != nil {
		return err
	}
	i, err := strconv.ParseFloat(rule, 64)
	if err != nil {
		return err
	}
	if val < i {
		return Error(name, fmt.Sprintf("must be greater than or equal to %f", i))
	}
	return nil
}

func getNumber(value any) (float64, error) {
	switch t := value.(type) {
	case int:
		{
			return float64(t), nil
		}
	case int16:
		{
			return float64(t), nil
		}
	case int32:
		{
			return float64(t), nil
		}
	case int64:
		{
			return float64(t), nil
		}
	case int8:
		{
			return float64(t), nil
		}
	case float32:
		{
			return float64(t), nil
		}
	case float64:
		{
			return t, nil
		}
	}
	return 0, fmt.Errorf("expected number but recieved %T", value)
}
