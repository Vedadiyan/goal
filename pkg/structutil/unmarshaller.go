package structutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	_ "unsafe"
)

var (
	_unmarshallers map[int]func(data map[string]any, field reflect.StructField, reflect reflect.Value) error
)

func init() {
	_unmarshallers = make(map[int]func(data map[string]any, field reflect.StructField, reflect reflect.Value) error)
	_unmarshallers[int(reflect.Float64)] = UnmarshalDouble
	_unmarshallers[int(reflect.Float32)] = UnmarshalFloat
	_unmarshallers[int(reflect.Int64)] = UnmarshalInt64
	_unmarshallers[int(reflect.Uint64)] = UnmarshalUInt64
	_unmarshallers[int(reflect.Int32)] = UnmarshalInt32
	_unmarshallers[int(reflect.Uint32)] = UnmarshalUInt32
	_unmarshallers[int(reflect.Int)] = UnmarshalInt
	_unmarshallers[int(reflect.Uint)] = UnmarshalUInt
	_unmarshallers[int(reflect.Bool)] = UnmarshalBool
	_unmarshallers[int(reflect.String)] = UnmarshalString
	_unmarshallers[int(reflect.Struct)] = UnmarshalMessage
	_unmarshallers[int(reflect.Int8)] = UnmarshalByte

	_unmarshallers[int(reflect.Float64)*100] = UnmarshalDoubleList
	_unmarshallers[int(reflect.Float32)*100] = UnmarshalFloatList
	_unmarshallers[int(reflect.Int64)*100] = UnmarshalInt64List
	_unmarshallers[int(reflect.Uint64)*100] = UnmarshalUInt64List
	_unmarshallers[int(reflect.Int32)*100] = UnmarshalInt32List
	_unmarshallers[int(reflect.Uint32)*100] = UnmarshalUInt32List
	_unmarshallers[int(reflect.Int)*100] = UnmarshalIntList
	_unmarshallers[int(reflect.Uint)*100] = UnmarshalUIntList
	_unmarshallers[int(reflect.Bool)*100] = UnmarshalBoolList
	_unmarshallers[int(reflect.String)*100] = UnmarshalStringList
	_unmarshallers[int(reflect.Struct)*100] = UnmarshalMessageList
	_unmarshallers[int(reflect.Int8)*100] = UnmarshalByteList
	_unmarshallers[int(reflect.Map)] = UnmarshalMessageMap
	_unmarshallers[int(reflect.Map)*100] = UnmarshalMessageMapList
	_unmarshallers[int(reflect.Slice)*100] = UnmarshalSlice
}

func Protect(err *error) {
	if r := recover(); r != nil {
		if r, ok := r.(error); ok {
			*err = r
		}
		*err = fmt.Errorf("%v", r)
	}
}

func UnmarshalDouble(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(doubleValue))
	return nil
}

func UnmarshalDoubleList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]float64, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		doubleValue, err := strconv.ParseFloat(valueRaw, 64)
		if err != nil {
			return err
		}
		slice = append(slice, doubleValue)
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalFloat(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(floatValue))
	return nil
}

func UnmarshalFloatList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]float32, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		floatValue, err := strconv.ParseFloat(valueRaw, 32)
		if err != nil {
			return err
		}
		slice = append(slice, float32(floatValue))
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalInt64(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(int64Value))
	return nil
}

func UnmarshalInt64List(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]int64, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		int64Value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil {
			return err
		}
		slice = append(slice, int64Value)
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalUInt64(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(uInt64Value))
	return nil
}

func UnmarshalUInt64List(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]uint64, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		uInt64Value, err := strconv.ParseUint(valueRaw, 10, 64)
		if err != nil {
			return err
		}
		slice = append(slice, uInt64Value)
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalInt32(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(int32(int32Value)))
	return nil
}

func UnmarshalInt32List(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]int32, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		int32Value, err := strconv.ParseInt(valueRaw, 10, 32)
		if err != nil {
			return err
		}
		slice = append(slice, int32(int32Value))
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalInt(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(int(int32Value)))
	return nil
}

func UnmarshalIntList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]int, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		int32Value, err := strconv.ParseInt(valueRaw, 10, 32)
		if err != nil {
			return err
		}
		slice = append(slice, int(int32Value))
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalUInt32(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(uint32(uInt32Value)))
	return nil
}

func UnmarshalUInt32List(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]uint32, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		uInt32Value, err := strconv.ParseUint(valueRaw, 10, 32)
		if err != nil {
			return err
		}
		slice = append(slice, uint32(uInt32Value))
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalUInt(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	v.Set(reflect.ValueOf(uint(uInt32Value)))
	return nil
}

func UnmarshalUIntList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]uint, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		uInt32Value, err := strconv.ParseUint(valueRaw, 10, 32)
		if err != nil {
			return err
		}
		slice = append(slice, uint(uInt32Value))
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalBool(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	v.Set(reflect.ValueOf(strings.ToLower(valueRaw) == "true"))
	return nil
}

func UnmarshalBoolList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]bool, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		slice = append(slice, strings.ToLower(valueRaw) == "true")
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalByte(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	byte, ok := value.(byte)
	if !ok {
		return fmt.Errorf("expected []byte but found %T", value)
	}
	v.Set(reflect.ValueOf(byte))
	return nil
}

func UnmarshalByteList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]byte, 0)
	for _, item := range list {
		bytes, ok := item.(byte)
		if !ok {
			return fmt.Errorf("expected []byte but found %T", value)
		}
		slice = append(slice, bytes)
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalString(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw := fmt.Sprintf("%v", value)
	v.Set(reflect.ValueOf(valueRaw))
	return nil
}

func UnmarshalStringList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := make([]string, 0)
	for _, item := range list {
		valueRaw := fmt.Sprintf("%v", item)
		slice = append(slice, valueRaw)
	}
	v.Set(reflect.ValueOf(slice))
	return nil
}

func UnmarshalMessage(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	message := reflect.New(f.Type)
	return Unmarshal(valueRaw, message.Interface())
}

func UnmarshalMessageMap(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	message := reflect.MapOf(f.Type.Key(), f.Type.Elem())
	mapper := reflect.MakeMap(message)
	for key, value := range valueRaw {

		kind := GetKindRaw(f.Type.Elem().Kind())
		switch kind {
		case int(reflect.Struct):
			{
				val := reflect.New(f.Type.Elem())
				err := Unmarshal(value.(map[string]any), val.Interface())
				if err != nil {
					return err
				}
				mapper.SetMapIndex(reflect.ValueOf(key), val.Elem())

			}
		default:
			{
				mapper.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
			}
		}
	}
	v.Set(mapper)
	return nil
}

func UnmarshalMessageMapList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	valueRaw, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected object by found %T", value)
	}
	slice := reflect.MakeSlice(reflect.SliceOf(f.Type.Elem()), 0, 0)
	for _, item := range valueRaw {
		mapper := reflect.MakeMap(reflect.MapOf(f.Type.Elem().Key(), f.Type.Elem().Elem()))
		for key, value := range item.(map[string]any) {

			kind := GetKindRaw(f.Type.Elem().Kind())
			switch kind {
			case int(reflect.Struct):
				{
					val := reflect.New(f.Type.Elem())
					err := Unmarshal(value.(map[string]any), val.Interface())
					if err != nil {
						return err
					}
					mapper.SetMapIndex(reflect.ValueOf(key), val.Elem())

				}
			default:
				{
					mapper.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
				}
			}
		}
		slice = reflect.Append(slice, mapper)
	}

	v.Set(slice)
	return nil
}

func UnmarshalMessageList(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
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
	slice := reflect.MakeSlice(reflect.SliceOf(f.Type.Elem()), 0, 0)
	for _, item := range list {
		valueRaw, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object by found %T", value)
		}
		message := reflect.New(f.Type.Elem())
		err := Unmarshal(valueRaw, message.Interface())
		if err != nil {
			return err
		}
		slice = reflect.Append(slice, message.Elem())
	}
	v.Set(slice)
	return nil
}

func UnmarshalSlice(d map[string]any, f reflect.StructField, v reflect.Value) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	depth := 0
	original := f.Type.Elem()
	for original.Kind() == reflect.Slice {
		original = original.Elem()
		depth++
	}
	var recursiveFn func(in any, depth int) reflect.Value
	recursiveFn = func(in any, d int) reflect.Value {
		v1 := reflect.ValueOf(in)
		if v1.Kind() == reflect.Map {
			v := reflect.New(original)
			err := Unmarshal(in.(map[string]any), v.Interface())
			if err != nil {

			}
			x := v.Interface()
			_ = x
			return v.Elem()
		}
		if v1.Kind() != reflect.Slice {
			return v1
		}
		sliceT := reflect.SliceOf(original)
		for i := 0; i < depth-d; i++ {
			sliceT = reflect.SliceOf(sliceT)
		}
		test := fmt.Sprintf("%v", sliceT)
		_ = test
		slice := reflect.MakeSlice(sliceT, 0, 0)
		for i := 0; i < v1.Len(); i++ {
			value := v1.Index(i).Interface()
			next := recursiveFn(value, d+1).Interface()
			slice = reflect.Append(slice, reflect.ValueOf(next))
		}
		return slice
	}
	tracker := recursiveFn(value, 0)
	v.Set(tracker)
	return nil
}

func GetKind(f reflect.StructField) int {
	field := f.Type
	if field.Kind() == reflect.Slice {
		return int(field.Elem().Kind()) * 100
	}
	return int(field.Kind())
}

func GetKindRaw(f reflect.Kind) int {
	if f == reflect.Slice {
		return int(f) * 100
	}
	return int(f)
}

func GetFieldName(field reflect.StructField) string {
	name := field.Tag.Get("name")
	if len(name) == 0 {
		name = field.Name
	}
	return name
}

func Unmarshal(data map[string]any, message any) error {
	p := reflect.TypeOf(message).Elem()
	v := reflect.ValueOf(message).Elem()
	n := p.NumField()
	for i := 0; i < n; i++ {
		field := p.Field(i)
		kind := GetKind(field)
		err := _unmarshallers[kind](data, field, v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}
