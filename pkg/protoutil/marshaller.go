package protoutil

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Marshal(message proto.Message) (map[string]any, error) {
	mapper := make(map[string]any)
	err := marshallerNext(message, "", mapper)
	if err != nil {
		return nil, err
	}
	return mapper, nil
}
func marshallerSetValue(data any, name string, index int, len int, v map[string]any) {
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
func marshallerSetList(ls protoreflect.List, name string, v map[string]any) error {
	for i := 0; i < ls.Len(); i++ {
		value := ls.Get(i)
		_, ok := value.Interface().(protoreflect.Message)
		if !ok {
			marshallerSetValue(value.Interface(), name, i, ls.Len(), v)
			continue
		}
		err := marshallerSetObject(value.Message().Interface(), name, i, ls.Len(), v)
		if err != nil {
			return err
		}
	}
	return nil
}
func marshallerSetObject(reflect proto.Message, name string, index int, len int, v map[string]any) error {
	mapper := make(map[string]any)
	err := marshallerNext(reflect, "", mapper)
	if err != nil {
		return err
	}
	marshallerSetValue(mapper, name, index, len, v)
	return nil
}
func marshallerNext(message proto.Message, name string, v map[string]any) error {
	reflect := message.ProtoReflect()
	fields := reflect.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		name := Or(field.JSONName(), field.TextName())
		if field.IsList() {
			isNotNull := reflect.Has(field)
			if isNotNull {
				err := marshallerSetList(reflect.Get(field).List(), name, v)
				if err != nil {
					return err
				}
				continue
			}
			marshallerSetValue([]any{}, name, -1, -1, v)
			continue
		}
		switch field.Kind() {
		case protoreflect.MessageKind:
			{
				isNotNull := reflect.Has(field)
				if isNotNull {
					err := marshallerSetObject(reflect.Get(field).Message().Interface(), name, -1, -1, v)
					if err != nil {
						return err
					}
				}
			}
		default:
			{
				isNotNull := reflect.Has(field)
				if isNotNull {
					marshallerSetValue(reflect.Get(field).Interface(), name, -1, -1, v)
				}
			}
		}
	}
	return nil
}
