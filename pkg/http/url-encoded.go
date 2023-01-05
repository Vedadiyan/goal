package http

import "reflect"

func ToURLEncoded(data any) Query {
	if data == nil {
		return nil
	}
	query := make(Query)
	v := reflect.ValueOf(data)
	reflectedType := v.Type()
	len := reflectedType.Elem().NumField()
	for i := 0; i < len; i++ {
		name := reflectedType.Elem().Field(i).Tag.Get("json")
		if name == "" {
			name = reflectedType.Elem().Field(i).Name
		}
		if v.Elem().Field(i).Kind() != reflect.Pointer {
			value := v.Elem().Field(i).String()
			query[name] = value
			continue
		}
		value := v.Elem().Field(i).Elem().String()
		query[name] = value
	}
	return query
}
