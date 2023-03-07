package protoutil

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ToMap(data proto.Message) (map[string]any, error) {
	mapper := make(map[string]any)
	err := next(data, "", mapper)
	if err != nil {
		return nil, err
	}
	return mapper, nil
}
func setValue(data any, name string, index int, len int, v map[string]any) {
	if index < 0 {
		v[name] = data
		return
	}
	_, ok := v[name]
	if !ok {
		v[name] = make([]any, len)
	}
	v[name].([]any)[index] = data
}
func setList(ls protoreflect.List, name string, v map[string]any) error {
	for x := 0; x < ls.Len(); x++ {
		value := ls.Get(x)
		_, ok := value.Interface().(protoreflect.Message)
		if !ok {
			setValue(value.Interface(), name, x, ls.Len(), v)
			continue
		}
		err := setObject(value.Message().Interface(), name, x, ls.Len(), v)
		if err != nil {
			return err
		}
	}
	return nil
}
func setObject(reflect proto.Message, name string, index int, len int, v map[string]any) error {
	mapper := make(map[string]any)
	err := next(reflect, "", mapper)
	if err != nil {
		return err
	}
	setValue(mapper, name, index, len, v)
	return nil
}
func next(message proto.Message, name string, v map[string]any) error {
	reflect := message.ProtoReflect()
	fields := reflect.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		name := Or(field.JSONName(), field.TextName())
		if field.IsList() {
			isNotNull := reflect.Has(field)
			if isNotNull {
				err := setList(reflect.Get(field).List(), name, v)
				if err != nil {
					return err
				}
			}
			continue
		}
		switch field.Kind() {
		case protoreflect.MessageKind:
			{
				isNotNull := reflect.Has(field)
				if isNotNull {
					err := setObject(reflect.Get(field).Message().Interface(), name, -1, -1, v)
					if err != nil {
						return err
					}
				}
			}
		default:
			{
				isNotNull := reflect.Has(field)
				if isNotNull {
					setValue(reflect.Get(field).Interface(), name, -1, -1, v)
				}
			}
		}
	}
	return nil
}

func Or(params ...string) string {
	for i := 0; i < len(params); i++ {
		if params[i] != "" {
			return params[i]
		}
	}
	return ""
}
