package protoutil

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func FromMap[T any](data map[string]any) (T, error) {
	message := new(T)
	if _, ok := any(message).(proto.Message); !ok {
		return *message, nil
	}
	err := _next(data, nil, -1, any(message).(proto.Message))
	if err != nil {
		return *message, err
	}
	return *message, nil
}

func _set(data any, fd protoreflect.FieldDescriptor, message proto.Message) {
	if !fd.IsList() {
		message.ProtoReflect().Set(fd, protoreflect.ValueOf(data))
		return
	}
	ls := message.(protoreflect.List)
	ls.Append(protoreflect.ValueOf(data))
}

func _setList(data []any, fields protoreflect.FieldDescriptors, fd protoreflect.FieldDescriptor, name string, message proto.Message) error {
	f := _getField(fields, name)
	if f == nil {
		return nil
	}
	ref := message.ProtoReflect().Mutable(f).List()
	for i := 0; i < len(data); i++ {
		switch t2 := data[i].(type) {
		case map[string]any:
			{
				ls := ref.AppendMutable().Message()
				err := _next(t2, f, i, ls.Interface())
				if err != nil {
					return err
				}
				ref.Set(i, protoreflect.ValueOfMessage(ls))
			}
		default:
			{
				ref.Append(protoreflect.ValueOf(data[i]))
			}
		}
	}
	message.ProtoReflect().Set(f, protoreflect.ValueOfList(ref))
	return nil
}

func _getField(fields protoreflect.FieldDescriptors, name string) protoreflect.FieldDescriptor {
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		if Or(f.JSONName(), f.TextName()) == name {
			return f
		}
	}
	return nil
}

func _setObject(data map[string]any, fields protoreflect.FieldDescriptors, name string, message proto.Message) error {
	f := _getField(fields, name)
	if f == nil {
		return nil
	}
	ref := message.ProtoReflect().Mutable(f).Message()
	err := _next(data, f, -1, ref.Interface())
	if err != nil {
		return err
	}
	_set(ref, f, message)
	return nil
}

func _setOneOf(value any, field protoreflect.FieldDescriptor, name string, message proto.Message) error {
	f := field.ContainingOneof().Fields().ByName(protoreflect.Name(name))
	switch f.Kind() {
	case protoreflect.MessageKind:
		{
			err := _next(value.(map[string]any), f, -1, message.ProtoReflect().Mutable(f).Message().New().Interface())
			if err != nil {
				return err
			}
		}
	default:
		{
			message.ProtoReflect().Set(f, protoreflect.ValueOf(value))
		}
	}
	return nil
}

func _setValue(value any, fields protoreflect.FieldDescriptors, name string, message proto.Message) error {
	f := _getField(fields, name)
	if f == nil {
		return nil
	}
	if f.ContainingOneof() != nil {
		err := _setOneOf(value, f, name, message)
		if err != nil {
			return nil
		}
		return nil
	}
	message.ProtoReflect().Set(f, protoreflect.ValueOf(value))
	return nil
}

func _next(data map[string]any, fd protoreflect.FieldDescriptor, index int, message proto.Message) error {
	fields := message.ProtoReflect().Descriptor().Fields()
	for key, value := range data {
		switch t := value.(type) {
		case map[string]any:
			{
				err := _setObject(t, fields, key, message)
				if err != nil {
					return err
				}
			}
		case []any:
			{
				err := _setList(t, fields, fd, key, message)
				if err != nil {
					return err
				}
			}
		default:
			{
				err := _setValue(value, fields, key, message)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

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
