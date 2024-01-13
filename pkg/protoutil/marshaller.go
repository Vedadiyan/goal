package protoutil

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	_marshallers map[Kind]func(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) error
)

func init() {
	_marshallers = make(map[Kind]func(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) error)
	_marshallers[DoubleKind] = MarshalDouble
	_marshallers[FloatKind] = MarshalFloat
	_marshallers[Int64Kind] = MarshalInt64
	_marshallers[Uint64Kind] = MarshalUInt64
	_marshallers[Int32Kind] = MarshalInt32
	_marshallers[Fixed64Kind] = MarshalUInt64
	_marshallers[Fixed32Kind] = MarshalUInt32
	_marshallers[BoolKind] = MarshalBool
	_marshallers[StringKind] = MarshalString
	_marshallers[GroupKind] = MarshalGroup
	_marshallers[MessageKind] = MarshalMessage
	_marshallers[BytesKind] = MarshalBytes
	_marshallers[Uint32Kind] = MarshalUInt32
	_marshallers[EnumKind] = MarshalEnum
	_marshallers[Sfixed32Kind] = MarshalInt32
	_marshallers[Sint32Kind] = MarshalInt32
	_marshallers[Sfixed64Kind] = MarshalInt64
	_marshallers[Sint64Kind] = MarshalInt64

	_marshallers[ListOfDoubleKind] = MarshalDoubleList
	_marshallers[ListOfFloatKind] = MarshalFloatList
	_marshallers[ListOfInt64Kind] = MarshalInt64List
	_marshallers[ListOfUint64Kind] = MarshalUInt64List
	_marshallers[ListOfInt32Kind] = MarshalInt32List
	_marshallers[ListOfFixed64Kind] = MarshalUInt64List
	_marshallers[ListOfFixed32Kind] = MarshalUInt32List
	_marshallers[ListOfBoolKind] = MarshalBoolList
	_marshallers[ListOfStringKind] = MarshalStringList
	_marshallers[ListOfGroupKind] = MarshalGroup
	_marshallers[ListOfMessageKind] = MarshalMessageList
	_marshallers[ListOfBytesKind] = MarshalBytesList
	_marshallers[ListOfUint32Kind] = MarshalUInt32List
	_marshallers[ListOfEnumKind] = MarshalEnumList
	_marshallers[ListOfSfixed32Kind] = MarshalInt32List
	_marshallers[ListOfSint32Kind] = MarshalInt32List
	_marshallers[ListOfSfixed64Kind] = MarshalInt64List
	_marshallers[ListOfSint64Kind] = MarshalInt64List
	_marshallers[MapKind] = MarshalMessageMap

	_marshallers[StructKind] = MarshalStruct
}

func MarshalDouble(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = value.Float()
	return nil
}

func MarshalDoubleList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]float64, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, value.Get(i).Float())
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalFloat(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = float32(value.Float())
	return nil
}

func MarshalFloatList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]float32, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, float32(value.Get(i).Float()))
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalInt64(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = value.Int()
	return nil
}

func MarshalInt64List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]int64, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, value.Get(i).Int())
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalUInt64(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = value.Uint()
	return nil
}

func MarshalUInt64List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]uint64, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, value.Get(i).Uint())
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalInt32(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = int32(value.Int())
	return nil
}

func MarshalInt32List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]int32, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, int32(value.Get(i).Int()))
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalUInt32(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = uint32(value.Uint())
	return nil
}

func MarshalUInt32List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]uint32, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, uint32(value.Get(i).Int()))
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalBool(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = value.Bool()
	return nil
}

func MarshalBoolList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]bool, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, value.Get(i).Bool())
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalGroup(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	panic("not implemented")
}

func MarshalEnum(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = value.Enum()
	return nil

}

func MarshalEnumList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]int, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, int(value.Get(i).Enum()))
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalBytes(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = value.Bytes()
	return nil
}

func MarshalBytesList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([][]byte, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, value.Get(i).Bytes())
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalString(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	data[GetFieldName(field)] = value.String()
	return nil
}

func MarshalStringList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		slice = append(slice, value.Get(i).String())
	}
	data[GetFieldName(field)] = slice
	return nil
}

func MarshalStruct(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := any(reflect.Get(field).Message().Interface()).(*structpb.Struct)
	data[GetFieldName(field)] = value.AsMap()
	return nil
}

func MarshalMessage(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field)
	message, err := Marshal(value.Message().Interface())
	if err != nil {
		return err
	}
	data[GetFieldName(field)] = message
	return nil
}

func MarshalMessageMap(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (_error error) {
	defer Protect(&_error)
	f := field.(protoreflect.FieldDescriptor)
	value := reflect.Get(field).Map()
	mapper := make(map[string]any)
	var err error
	value.Range(func(mk protoreflect.MapKey, v protoreflect.Value) bool {
		err = _marshallers[GetKind(f.MapValue())](mapper, mk, MapType{Map: value})
		return err == nil
	})
	data[GetFieldName(field)] = mapper
	if err != nil {
		return err
	}
	return nil
}

func MarshalMessageList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value := reflect.Get(field).List()
	slice := make([]map[string]any, value.Len())
	for i := 0; i < value.Len(); i++ {
		value, err := Marshal(value.Get(i).Message().Interface())
		if err != nil {
			return err
		}
		slice = append(slice, value)
	}
	data[GetFieldName(field)] = slice
	return nil
}

func Marshal(p proto.Message) (map[string]any, error) {
	fields := p.ProtoReflect().Descriptor().Fields()
	mapper := make(map[string]any)
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		err := _marshallers[GetKind(field)](mapper, field, MessageType{Message: p.ProtoReflect()})
		if err != nil {
			return nil, err
		}
	}
	return mapper, nil
}
