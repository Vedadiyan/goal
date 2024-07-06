package protoutil

import (
	"fmt"
	"strconv"
	"strings"
	_ "unsafe"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

type (
	Kind                int
	FieldDescriptorKind any
	ProtobufType        interface {
		Get(FieldDescriptorKind) protoreflect.Value
		Mutable(FieldDescriptorKind) protoreflect.Value
		Set(FieldDescriptorKind, protoreflect.Value)
		HasValue(f FieldDescriptorKind) bool
	}
	MapType struct {
		Map protoreflect.Map
	}
	MessageType struct {
		Message protoreflect.Message
	}
)

const (
	DoubleKind         Kind = 1
	FloatKind          Kind = 2
	Int64Kind          Kind = 3
	Uint64Kind         Kind = 4
	Int32Kind          Kind = 5
	Fixed64Kind        Kind = 6
	Fixed32Kind        Kind = 7
	BoolKind           Kind = 8
	StringKind         Kind = 9
	GroupKind          Kind = 10
	MessageKind        Kind = 11
	BytesKind          Kind = 12
	Uint32Kind         Kind = 13
	EnumKind           Kind = 14
	Sfixed32Kind       Kind = 15
	Sint32Kind         Kind = 16
	Sfixed64Kind       Kind = 17
	Sint64Kind         Kind = 18
	ListOfDoubleKind   Kind = 101
	ListOfFloatKind    Kind = 102
	ListOfInt64Kind    Kind = 103
	ListOfUint64Kind   Kind = 104
	ListOfInt32Kind    Kind = 105
	ListOfFixed64Kind  Kind = 106
	ListOfFixed32Kind  Kind = 107
	ListOfBoolKind     Kind = 108
	ListOfStringKind   Kind = 109
	ListOfGroupKind    Kind = 110
	ListOfMessageKind  Kind = 111
	ListOfBytesKind    Kind = 112
	ListOfUint32Kind   Kind = 113
	ListOfEnumKind     Kind = 114
	ListOfSfixed32Kind Kind = 115
	ListOfSint32Kind   Kind = 116
	ListOfSfixed64Kind Kind = 117
	ListOfSint64Kind   Kind = 118
	MapKind            Kind = 119
	StructKind         Kind = 120
)

var (
	_unmarshallers map[Kind]func(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) error
)

func (mapMessage MapType) Get(f FieldDescriptorKind) protoreflect.Value {
	return mapMessage.Map.Get(f.(protoreflect.MapKey))
}
func (mapMessage MapType) Mutable(f FieldDescriptorKind) protoreflect.Value {
	return mapMessage.Map.Mutable(f.(protoreflect.MapKey))
}
func (mapMessage MapType) Set(f FieldDescriptorKind, v protoreflect.Value) {
	mapMessage.Map.Set(f.(protoreflect.MapKey), v)
}
func (mapMessage MapType) HasValue(f FieldDescriptorKind) bool {
	return mapMessage.Map.Has(f.(protoreflect.MapKey))
}

func (messageType MessageType) Get(f FieldDescriptorKind) protoreflect.Value {
	return messageType.Message.Get(f.(protoreflect.FieldDescriptor))
}
func (messageType MessageType) Mutable(f FieldDescriptorKind) protoreflect.Value {
	return messageType.Message.Mutable(f.(protoreflect.FieldDescriptor))
}
func (messageType MessageType) Set(f FieldDescriptorKind, v protoreflect.Value) {
	messageType.Message.Set(f.(protoreflect.FieldDescriptor), v)
}
func (messageType MessageType) HasValue(f FieldDescriptorKind) bool {
	return messageType.Message.Has(f.(protoreflect.FieldDescriptor))
}

func init() {
	_unmarshallers = make(map[Kind]func(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) error)
	_unmarshallers[DoubleKind] = UnmarshalDouble
	_unmarshallers[FloatKind] = UnmarshalFloat
	_unmarshallers[Int64Kind] = UnmarshalInt64
	_unmarshallers[Uint64Kind] = UnmarshalUInt64
	_unmarshallers[Int32Kind] = UnmarshalInt32
	_unmarshallers[Fixed64Kind] = UnmarshalUInt64
	_unmarshallers[Fixed32Kind] = UnmarshalUInt32
	_unmarshallers[BoolKind] = UnmarshalBool
	_unmarshallers[StringKind] = UnmarshalString
	_unmarshallers[GroupKind] = UnmarshalGroup
	_unmarshallers[MessageKind] = UnmarshalMessage
	_unmarshallers[BytesKind] = UnmarshalBytes
	_unmarshallers[Uint32Kind] = UnmarshalUInt32
	_unmarshallers[EnumKind] = UnmarshalEnum
	_unmarshallers[Sfixed32Kind] = UnmarshalInt32
	_unmarshallers[Sint32Kind] = UnmarshalInt32
	_unmarshallers[Sfixed64Kind] = UnmarshalInt64
	_unmarshallers[Sint64Kind] = UnmarshalInt64

	_unmarshallers[ListOfDoubleKind] = UnmarshalDoubleList
	_unmarshallers[ListOfFloatKind] = UnmarshalFloatList
	_unmarshallers[ListOfInt64Kind] = UnmarshalInt64List
	_unmarshallers[ListOfUint64Kind] = UnmarshalUInt64List
	_unmarshallers[ListOfInt32Kind] = UnmarshalInt32List
	_unmarshallers[ListOfFixed64Kind] = UnmarshalUInt64List
	_unmarshallers[ListOfFixed32Kind] = UnmarshalUInt32List
	_unmarshallers[ListOfBoolKind] = UnmarshalBoolList
	_unmarshallers[ListOfStringKind] = UnmarshalStringList
	_unmarshallers[ListOfGroupKind] = UnmarshalGroup
	_unmarshallers[ListOfMessageKind] = UnmarshalMessageList
	_unmarshallers[ListOfBytesKind] = UnmarshalBytesList
	_unmarshallers[ListOfUint32Kind] = UnmarshalUInt32List
	_unmarshallers[ListOfEnumKind] = UnmarshalEnumList
	_unmarshallers[ListOfSfixed32Kind] = UnmarshalInt32List
	_unmarshallers[ListOfSint32Kind] = UnmarshalInt32List
	_unmarshallers[ListOfSfixed64Kind] = UnmarshalInt64List
	_unmarshallers[ListOfSint64Kind] = UnmarshalInt64List
	_unmarshallers[MapKind] = UnmarshalMessageMap

	_unmarshallers[StructKind] = UnmarshalStruct
}

func Protect(err *error) {
	if r := recover(); r != nil {
		if r, ok := r.(error); ok {
			*err = r
		}
		*err = fmt.Errorf("%v", r)
	}
}

func UnmarshalDouble(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	doubleValue, err := strconv.ParseFloat(valueRaw, 64)
	if err != nil {
		return err
	}
	reflect.Set(field, protoreflect.ValueOf(doubleValue))
	return nil
}

func UnmarshalDoubleList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		doubleValue, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return err
		}
		v.Append(protoreflect.ValueOf(doubleValue))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalFloat(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	floatValue, err := strconv.ParseFloat(valueRaw, 32)
	if err != nil {
		return err
	}
	reflect.Set(field, protoreflect.ValueOf(float32(floatValue)))
	return nil
}

func UnmarshalFloatList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		floatValue, err := strconv.ParseFloat(valueRaw, 32)
		if err != nil {
			return err
		}
		v.Append(protoreflect.ValueOf(floatValue))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalInt64(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	int64Value, err := strconv.ParseInt(valueRaw, 10, 64)
	if err != nil {
		return err
	}
	reflect.Set(field, protoreflect.ValueOf(int64Value))
	return nil
}

func UnmarshalInt64List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		int64Value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			return err
		}
		v.Append(protoreflect.ValueOf(int64Value))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalUInt64(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	uInt64Value, err := strconv.ParseUint(valueRaw, 10, 64)
	if err != nil {
		return err
	}
	reflect.Set(field, protoreflect.ValueOf(uInt64Value))
	return nil
}

func UnmarshalUInt64List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		uInt64Value, err := strconv.ParseUint(valueRaw, 10, 64)
		if err != nil {
			return err
		}
		v.Append(protoreflect.ValueOf(uInt64Value))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalInt32(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	int32Value, err := strconv.ParseInt(valueRaw, 10, 32)
	if err != nil {
		return err
	}
	reflect.Set(field, protoreflect.ValueOf(int32(int32Value)))
	return nil
}

func UnmarshalInt32List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		int32Value, err := strconv.ParseInt(valueRaw, 10, 32)
		if err != nil {
			return err
		}
		v.Append(protoreflect.ValueOf(int32(int32Value)))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalUInt32(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	uInt32Value, err := strconv.ParseUint(valueRaw, 10, 32)
	if err != nil {
		return err
	}
	reflect.Set(field, protoreflect.ValueOf(uint32(uInt32Value)))
	return nil
}

func UnmarshalUInt32List(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		uInt32Value, err := strconv.ParseUint(valueRaw, 10, 32)
		if err != nil {
			return err
		}
		v.Append(protoreflect.ValueOf(uint32(uInt32Value)))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalBool(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	reflect.Set(field, protoreflect.ValueOf(strings.ToLower(valueRaw) == "true"))
	return nil
}

func UnmarshalBoolList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		v.Append(protoreflect.ValueOf(strings.ToLower(valueRaw) == "true"))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalGroup(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	panic("not implemented")
}

func UnmarshalEnum(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string by found %T", value)
	}
	enum := field.(protoreflect.FieldDescriptor).Enum().Values().ByName(protoreflect.Name(str))
	if enum == nil {
		return nil
	}
	reflect.Set(field, protoreflect.ValueOf(enum.Number()))
	return nil
}

func UnmarshalEnumList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		str, ok := item.(string)
		if !ok {
			return fmt.Errorf("expected string by found %T", value)
		}
		enum := field.(protoreflect.FieldDescriptor).Enum().Values().ByName(protoreflect.Name(str))
		if enum == nil {
			return nil
		}
		v.Append(protoreflect.ValueOf(enum.Number()))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalBytes(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte but found %T", value)
	}
	reflect.Set(field, protoreflect.ValueOf(bytes))
	return nil
}

func UnmarshalBytesList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		bytes, ok := item.([]byte)
		if !ok {
			return fmt.Errorf("expected []byte but found %T", value)
		}
		v.Append(protoreflect.ValueOf(bytes))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalString(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	reflect.Set(field, protoreflect.ValueOf(valueRaw))
	return nil
}

func UnmarshalStringList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		v.Append(protoreflect.ValueOf(valueRaw))
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func UnmarshalStruct(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	_value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	data, ok = _value.((map[string]any))
	if !ok {
		return fmt.Errorf("expected object by found %T", _value)
	}
	value, err := structpb.NewStruct(data)
	if err != nil {
		return nil
	}
	reflect.Set(field, protoreflect.ValueOfMessage(value.ProtoReflect()))
	return nil
}

func UnmarshalMessage(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("expected object by found %T", value)
	}
	message := reflect.Mutable(field).Message().Interface()
	return Unmarshal(valueRaw, message)
}

func UnmarshalMessageMap(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("expected object by found %T", value)
	}
	message := reflect.Mutable(field).Map()
	for key, value := range valueRaw {
		field := field.(protoreflect.FieldDescriptor)
		kind := GetKind(field.MapValue())
		switch kind {
		case MessageKind:
			{
				val := message.NewValue()
				inst := val.Message().Interface()
				err := Unmarshal(value.(map[string]any), inst)
				if err != nil {
					return err
				}
				message.Set(protoreflect.MapKey(protoreflect.ValueOf(key)), val)
			}
		default:
			{
				err := _unmarshallers[GetKind(field.MapValue())](valueRaw, protoreflect.MapKey(protoreflect.ValueOf(key)), MapType{Map: message})
				if err != nil {
					return err
				}
			}
		}
	}
	reflect.Set(field, protoreflect.ValueOf(message))
	return nil
}

func UnmarshalMessageList(data map[string]any, field FieldDescriptorKind, reflect ProtobufType) (error error) {
	defer Protect(&error)
	value, ok := data[GetFieldName(field)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	v := reflect.Mutable(field).List()
	for _, item := range list {
		if item == nil {
			v.Append(protoreflect.ValueOf(reflect.Get(field).List().NewElement().Message().Type().Zero()))
			continue
		}
		valueRaw, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object by found %T", value)
		}
		elem := reflect.Get(field).List().NewElement()
		message := elem.Message().Interface()
		err := Unmarshal(valueRaw, message)
		if err != nil {
			return err
		}
		v.Append(elem)
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func GetKind(f FieldDescriptorKind) Kind {
	field := f.(protoreflect.FieldDescriptor)
	if field.IsList() {
		return Kind(field.Kind() + 100)
	}
	if field.IsMap() {
		return MapKind
	}
	message := field.Message()
	if message != nil {
		switch message.FullName() {
		case "google.protobuf.Struct":
			{
				return StructKind
			}
		}
	}
	return Kind(field.Kind())
}

func GetFieldName(field FieldDescriptorKind) string {
	switch field := field.(type) {
	case protoreflect.MapKey:
		{
			return field.Value().String()
		}
	case protoreflect.FieldDescriptor:
		{
			if len(field.JSONName()) != 0 {
				return field.JSONName()
			}
			return field.TextName()
		}
	default:
		{
			panic("")
		}
	}
}

func Unmarshal(data map[string]any, message any) error {
	p := message.(protoreflect.ProtoMessage)
	reflect := p.ProtoReflect()
	descriptor := reflect.Descriptor()
	fields := descriptor.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		kind := GetKind(field)
		err := _unmarshallers[kind](data, field, MessageType{Message: reflect})
		if err != nil {
			return err
		}
	}
	return nil
}
