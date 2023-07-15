package protoutil

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

func Unmarshal(data any, message proto.Message) error {
	if structPb, ok := message.(*structpb.Struct); ok {
		switch t := data.(type) {
		case map[string]any:
			{
				res, err := structpb.NewStruct(t)
				if err != nil {
					return err
				}
				*structPb = *res
			}
		default:
			{
				res, err := structpb.NewStruct(map[string]any{"data": t})
				if err != nil {
					return err
				}
				*structPb = *res
			}
		}
		return nil
	}
	return unmarshallerNext(data, nil, -1, message)
}
func unmarshallerSet(data any, fd protoreflect.FieldDescriptor, message proto.Message) {
	if !fd.IsList() {
		message.ProtoReflect().Set(fd, protoreflect.ValueOf(castType(data)))
		return
	}
	ls := message.(protoreflect.List)
	ls.Append(protoreflect.ValueOf(castType(data)))
}
func unmarshallerSetList(data []any, fields protoreflect.FieldDescriptors, fd protoreflect.FieldDescriptor, name string, message proto.Message) error {
	f := unmarshallerGetField(fields, name)
	if f == nil {
		return nil
	}
	ref := message.ProtoReflect().Mutable(f).List()
	for i := 0; i < len(data); i++ {
		switch t2 := data[i].(type) {
		case map[string]any:
			{
				ls := ref.AppendMutable().Message()
				err := unmarshallerNext(t2, f, i, ls.Interface())
				if err != nil {
					return err
				}
				ref.Set(i, protoreflect.ValueOfMessage(ls))
			}
		default:
			{
				ref.Append(protoreflect.ValueOf(castType(data[i])))
			}
		}
	}
	message.ProtoReflect().Set(f, protoreflect.ValueOfList(ref))
	return nil
}
func unmarshallerGetField(fields protoreflect.FieldDescriptors, name string) protoreflect.FieldDescriptor {
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		if Or(f.JSONName(), f.TextName()) == name {
			return f
		}
	}
	return nil
}
func unmarshallerSetObject(data map[string]any, fields protoreflect.FieldDescriptors, name string, message proto.Message) error {
	f := unmarshallerGetField(fields, name)
	if f == nil {
		return nil
	}
	if f.IsMap() {
		ref := message.ProtoReflect().Mutable(f).Map()
		for key, value := range data {
			switch t := value.(type) {
			case map[string]any:
				{
					// WARNING: UNTESTED CODE
					value := ref.NewValue().Message().Interface()
					err := unmarshallerNext(t, f, -1, value)
					if err != nil {
						return err
					}
					ref.Set(protoreflect.ValueOf(key).MapKey(), protoreflect.ValueOf(value))
				}
			default:
				{
					ref.Set(protoreflect.ValueOf(key).MapKey(), protoreflect.ValueOf(value))
				}
			}
		}
		message.ProtoReflect().Set(f, protoreflect.ValueOfMap(ref))
		return nil
	}
	ref := message.ProtoReflect().Mutable(f).Message()
	err := unmarshallerNext(data, f, -1, ref.Interface())
	if err != nil {
		return err
	}
	unmarshallerSet(ref, f, message)
	return nil
}
func unmarshallerSetOneOf(value any, field protoreflect.FieldDescriptor, name string, message proto.Message) error {
	f := field.ContainingOneof().Fields().ByName(protoreflect.Name(name))
	switch f.Kind() {
	case protoreflect.MessageKind:
		{
			err := unmarshallerNext(value.(map[string]any), f, -1, message.ProtoReflect().Mutable(f).Message().New().Interface())
			if err != nil {
				return err
			}
		}
	default:
		{
			message.ProtoReflect().Set(f, protoreflect.ValueOf(castType(value)))
		}
	}
	return nil
}
func unmarshallerSetValue(value any, fields protoreflect.FieldDescriptors, name string, message proto.Message) error {
	f := unmarshallerGetField(fields, name)
	if f == nil {
		return nil
	}
	if f.ContainingOneof() != nil {
		err := unmarshallerSetOneOf(value, f, name, message)
		if err != nil {
			return nil
		}
		return nil
	}
	message.ProtoReflect().Set(f, protoreflect.ValueOf(castType(value)))
	return nil
}
func unmarshallerNext(data any, fd protoreflect.FieldDescriptor, index int, message proto.Message) error {
	fields := message.ProtoReflect().Descriptor().Fields()
	switch t := data.(type) {
	case map[string]any:
		{
			for key, value := range t {
				switch t := value.(type) {
				case map[string]any:
					{
						err := unmarshallerSetObject(t, fields, key, message)
						if err != nil {
							return err
						}
					}
				case []any:
					{
						err := unmarshallerSetList(t, fields, fd, key, message)
						if err != nil {
							return err
						}
					}
				default:
					{
						if value == nil {
							continue
						}
						err := unmarshallerSetValue(value, fields, key, message)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	case []any:
		{
			for _, value := range t {
				err := unmarshallerNext(value, fd, index, message)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func castType(value any) any {
	switch t := value.(type) {
	case int:
		{
			return int32(t)
		}
	default:
		{
			return t
		}
	}
}
