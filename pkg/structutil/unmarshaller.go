package structutil

import (
	"fmt"
	"reflect"
	_ "unsafe"
)

var (
	_unmarshallers map[int]func(data map[string]any, field reflect.StructField, reflect reflect.Value, pointerDepth int) error
)

func init() {
	_unmarshallers = make(map[int]func(data map[string]any, field reflect.StructField, reflect reflect.Value, pointerDepth int) error)
	_unmarshallers[int(reflect.Int8)] = UnmarshallSimple[int8]
	_unmarshallers[int(reflect.Int16)] = UnmarshallSimple[int16]
	_unmarshallers[int(reflect.Int32)] = UnmarshallSimple[int32]
	_unmarshallers[int(reflect.Int)] = UnmarshallSimple[int]
	_unmarshallers[int(reflect.Int64)] = UnmarshallSimple[int64]

	_unmarshallers[int(reflect.Uint8)] = UnmarshallSimple[uint8]
	_unmarshallers[int(reflect.Uint16)] = UnmarshallSimple[uint16]
	_unmarshallers[int(reflect.Uint32)] = UnmarshallSimple[uint32]
	_unmarshallers[int(reflect.Uint)] = UnmarshallSimple[uint]
	_unmarshallers[int(reflect.Uint64)] = UnmarshallSimple[uint64]

	_unmarshallers[int(reflect.Float32)] = UnmarshallSimple[float32]
	_unmarshallers[int(reflect.Float64)] = UnmarshallSimple[float64]
	_unmarshallers[int(reflect.Bool)] = UnmarshallSimple[bool]
	_unmarshallers[int(reflect.String)] = UnmarshallSimple[string]

	_unmarshallers[int(reflect.Int8)*100] = UnmarshallSimpleSlice[int8]
	_unmarshallers[int(reflect.Int16)*100] = UnmarshallSimpleSlice[int16]
	_unmarshallers[int(reflect.Int32)*100] = UnmarshallSimpleSlice[int32]
	_unmarshallers[int(reflect.Int)*100] = UnmarshallSimpleSlice[int]
	_unmarshallers[int(reflect.Int64)*100] = UnmarshallSimpleSlice[int64]

	_unmarshallers[int(reflect.Uint8)*100] = UnmarshallSimpleSlice[uint8]
	_unmarshallers[int(reflect.Uint16)*100] = UnmarshallSimpleSlice[uint16]
	_unmarshallers[int(reflect.Uint32)*100] = UnmarshallSimpleSlice[uint32]
	_unmarshallers[int(reflect.Uint)*100] = UnmarshallSimpleSlice[uint]
	_unmarshallers[int(reflect.Uint64)*100] = UnmarshallSimpleSlice[uint64]

	_unmarshallers[int(reflect.Float32)*100] = UnmarshallSimpleSlice[float32]
	_unmarshallers[int(reflect.Float64)*100] = UnmarshallSimpleSlice[float64]
	_unmarshallers[int(reflect.Bool)*100] = UnmarshallSimpleSlice[bool]
	_unmarshallers[int(reflect.String)*100] = UnmarshallSimpleSlice[string]

	_unmarshallers[int(reflect.Struct)] = UnmarshalMessage
	_unmarshallers[int(reflect.Pointer)] = UnmarshalPointer
	_unmarshallers[int(reflect.Pointer)*100] = UnmarshalPointerSlice

	_unmarshallers[int(reflect.Struct)*100] = UnmarshalMessageList

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

func UnmarshallSimple[T any](d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	kind := reflect.ValueOf(new(T)).Elem().Kind()
	output, err := _convertors[kind](value)
	if err != nil {
		return err
	}
	Set(output.(T), v, pointerDepth)
	return nil
}

func UnmarshallSimpleSlice[T any](d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
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
	kind := reflect.ValueOf(new(T)).Elem().Kind()
	slice := make([]T, 0)
	for _, item := range list {
		value, err := _convertors[kind](item)
		if err != nil {
			return err
		}
		slice = append(slice, value.(T))
	}
	Set(slice, v, pointerDepth)
	return nil
}

func Set(value any, v reflect.Value, refrenceDepth int) {
	if refrenceDepth != 0 {
		pointerType := reflect.TypeOf(value)
		pointerValue := reflect.New(pointerType)
		pointerValue.Elem().Set(reflect.ValueOf(reflect.ValueOf(&value).Elem().Interface()))
		for i := 1; i < refrenceDepth; i++ {
			pointerType = reflect.PointerTo(pointerType)
			temp := reflect.New(pointerType)
			temp.Elem().Set(pointerValue)
			pointerValue = temp
		}
		v.Set(pointerValue)
		return
	}
	v.Set(reflect.ValueOf(value))
}

func UnmarshalMessage(d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	if f.Type == reflect.TypeOf(value) {
		v.Set(reflect.ValueOf(value))
		return nil
	}
	valueRaw, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("expected object by found %T", value)
	}
	message := reflect.New(f.Type)
	messageInterface := message.Interface()
	err := Unmarshal(valueRaw, messageInterface)
	if err != nil {
		return error
	}
	Set(message.Elem().Interface(), v, pointerDepth)
	return nil
}

func UnmarshalMessageList(d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
	defer Protect(&error)
	value, ok := d[GetFieldName(f)]
	if !ok {
		return nil
	}
	if value == nil {
		return nil
	}
	if f.Type == reflect.TypeOf(value) {
		v.Set(reflect.ValueOf(value))
		return nil
	}
	list, ok := value.([]any)
	if !ok {
		return fmt.Errorf("expected list by found %T", value)
	}
	slice := reflect.MakeSlice(reflect.SliceOf(f.Type.Elem()), 0, 0)
	for _, item := range list {
		if f.Type == reflect.TypeOf(item) {
			slice = reflect.Append(slice, reflect.ValueOf(item))
		}
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

func UnmarshalMessageMap(d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
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
	Set(mapper.Interface(), v, pointerDepth)
	return nil
}

func UnmarshalMessageMapList(d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
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

func UnmarshalSlice(d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
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
			if original.Kind() == reflect.Map {
				return v1
			}
			if original.Kind() == reflect.Interface {
				return v1
			}
			_ = Unmarshal(in.(map[string]any), v.Interface())
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
	Set(tracker.Interface(), v, pointerDepth)
	return nil
}

func UnmarshalPointer(d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
	defer Protect(&error)
	f.Type = f.Type.Elem()
	return _unmarshallers[GetKindRaw(f.Type.Kind())](d, f, v, pointerDepth+1)
}

func UnmarshalPointerSlice(d map[string]any, f reflect.StructField, v reflect.Value, pointerDepth int) (error error) {
	//defer Protect(&error)
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
	ptrDepth := 0
	original2 := f.Type.Elem()
	for original2.Kind() == reflect.Pointer {
		original2 = original2.Elem()
		ptrDepth++
	}
	var recursiveFn func(in any, depth int) reflect.Value
	recursiveFn = func(in any, d int) reflect.Value {
		v1 := reflect.ValueOf(in)
		if v1.Kind() == reflect.Map {
			v := reflect.New(original)
			if original.Kind() == reflect.Map {
				return v1
			}
			if original.Kind() == reflect.Interface {
				return v1
			}
			intfc := v
			if ptrDepth > 0 {
				for i := 0; i < ptrDepth; i++ {
					intfc = v.Elem()
					intfc.Set(reflect.New(original2))
				}
			}
			xx := intfc.Interface()
			_ = xx
			_ = Unmarshal(in.(map[string]any), intfc.Interface())
			return intfc.Elem()
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
			v := reflect.New(original)
			Set(next, v.Elem(), ptrDepth)
			slice = reflect.Append(slice, v.Elem())
		}
		return slice
	}
	tracker := recursiveFn(value, 0)
	Set(tracker.Interface(), v, pointerDepth)
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
		err := _unmarshallers[kind](data, field, v.Field(i), 0)
		if err != nil {
			return err
		}
	}
	return nil
}
