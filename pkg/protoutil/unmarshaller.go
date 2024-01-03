package protoutil

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

type (
	Kind int
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
	_handlers map[Kind]func(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) error
)

func init() {
	_handlers = make(map[Kind]func(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) error)
	_handlers[DoubleKind] = Double
	_handlers[FloatKind] = Float
	_handlers[Int64Kind] = Int64
	_handlers[Uint64Kind] = UInt64
	_handlers[Int32Kind] = Int32
	_handlers[Fixed64Kind] = UInt64
	_handlers[Fixed32Kind] = UInt32
	_handlers[BoolKind] = Bool
	_handlers[StringKind] = String
	_handlers[GroupKind] = Group
	_handlers[MessageKind] = Message
	_handlers[BytesKind] = Bytes
	_handlers[Uint32Kind] = UInt32
	_handlers[EnumKind] = Enum
	_handlers[Sfixed32Kind] = Int32
	_handlers[Sint32Kind] = Int32
	_handlers[Sfixed64Kind] = Int64
	_handlers[Sint64Kind] = Int64

	_handlers[ListOfDoubleKind] = DoubleList
	_handlers[ListOfFloatKind] = FloatList
	_handlers[ListOfInt64Kind] = Int64List
	_handlers[ListOfUint64Kind] = UInt64List
	_handlers[ListOfInt32Kind] = Int32List
	_handlers[ListOfFixed64Kind] = UInt64List
	_handlers[ListOfFixed32Kind] = UInt32List
	_handlers[ListOfBoolKind] = BoolList
	_handlers[ListOfStringKind] = StringList
	_handlers[ListOfGroupKind] = Group
	_handlers[ListOfMessageKind] = MessageList
	_handlers[ListOfBytesKind] = BytesList
	_handlers[ListOfUint32Kind] = UInt32List
	_handlers[ListOfEnumKind] = Enum
	_handlers[ListOfSfixed32Kind] = Int32List
	_handlers[ListOfSint32Kind] = Int32List
	_handlers[ListOfSfixed64Kind] = Int64List
	_handlers[ListOfSint64Kind] = Int64List
	_handlers[MapKind] = MessageMap

	_handlers[StructKind] = Struct
}

func Protect(err *error) {
	if r := recover(); r != nil {
		if r, ok := r.(error); ok {
			*err = r
		}
		*err = fmt.Errorf("%v", r)
	}
}

func Double(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func DoubleList(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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
	reflect.Set(field, protoreflect.ValueOf(list))
	return nil
}

func Float(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func FloatList(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func Int64(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func Int64List(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func UInt64(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func UInt64List(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func Int32(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func Int32List(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func UInt32(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func UInt32List(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func Bool(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func BoolList(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func Group(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
	panic("not implemented")
}

func Enum(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
	panic("not implemented")
}

func Bytes(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func BytesList(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func String(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func StringList(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func Struct(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
	defer Protect(&error)
	value, err := structpb.NewStruct(data)
	if err != nil {
		return nil
	}
	reflect.Set(field, protoreflect.ValueOfMessage(value.ProtoReflect()))
	return nil
}

func Message(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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

func MessageMap(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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
				message.Set(protoreflect.MapKey(protoreflect.ValueOf(key)), protoreflect.ValueOf(value))
			}
		}
	}
	reflect.Set(field, protoreflect.ValueOf(message))
	return nil
}

func MessageList(data map[string]any, field protoreflect.FieldDescriptor, reflect protoreflect.Message) (error error) {
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
		valueRaw, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object by found %T", value)
		}
		elem := reflect.Get(field).List().NewElement()
		message := elem.Message().Interface()
		err := Unmarshal(valueRaw, message)
		if err != nil {
			return nil
		}
		v.Append(elem)
	}
	reflect.Set(field, protoreflect.ValueOf(v))
	return nil
}

func GetKind(field protoreflect.FieldDescriptor) Kind {
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

func GetFieldName(field protoreflect.FieldDescriptor) string {
	if len(field.JSONName()) != 0 {
		return field.JSONName()
	}
	return field.TextName()
}

func Unmarshal(data map[string]any, message any) error {
	p := message.(protoreflect.ProtoMessage)
	reflect := p.ProtoReflect()
	descriptor := reflect.Descriptor()
	fields := descriptor.Fields()
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		kind := GetKind(field)
		err := _handlers[kind](data, field, reflect)
		if err != nil {
			return err
		}
	}
	return nil
}
