package structutil

import (
	"encoding/json"
	"testing"
)

type (
	Simple struct {
		Number  int
		Text    string
		Boolean bool
	}
	SimplePointer struct {
		Number  *int
		Text    string
		Boolean bool
	}
	SimpleDoublePointer struct {
		Number  **int
		Text    string
		Boolean bool
	}
	SimpleAllTypes struct {
		Int     int
		Int32   int32
		Int64   int64
		Int16   int16
		Int8    int8
		Uint    uint
		Uint32  uint32
		Uint64  uint64
		Uint16  uint16
		Uint8   uint8
		Byte    byte
		Text    string
		Boolean bool
	}
	SimpleAllArrayTypes struct {
		Int     []int
		Int32   []int32
		Int64   []int64
		Int16   []int16
		Int8    []int8
		Uint    []uint
		Uint32  []uint32
		Uint64  []uint64
		Uint16  []uint16
		Uint8   []uint8
		Byte    []uint8
		Text    []string
		Boolean []bool
	}
	SimpleArray struct {
		Number  []int
		Text    string
		Boolean bool
	}
	ArrayStruct struct {
		Number  []Simple
		Text    string
		Boolean bool
	}
	ArrayMap struct {
		Number  []map[string]any
		Text    string
		Boolean bool
	}
	ArrayMap2d struct {
		Number  [][]map[string]any
		Text    string
		Boolean bool
	}
	ArrayStructMulti struct {
		Number  [][]Simple
		Text    string
		Boolean bool
	}
	Array2D struct {
		Number  [][]int
		Text    string
		Boolean bool
	}
	ArrayMultiD struct {
		Number  [][][]int
		Text    string
		Boolean bool
	}
	SimpleMap struct {
		Number  map[string]int
		Text    string
		Boolean bool
	}
	MapStruct struct {
		Number  map[string]Simple
		Text    string
		Boolean bool
	}
)

func TestSimple(t *testing.T) {
	data := map[string]any{
		"Number":  10,
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(Simple)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestSimplePointer(t *testing.T) {
	data := map[string]any{
		"Number":  10,
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(SimplePointer)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestSimpleDoublePointer(t *testing.T) {
	data := map[string]any{
		"Number":  10,
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(SimpleDoublePointer)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestSimpleAllTypes(t *testing.T) {
	simpleAllTypes := SimpleAllTypes{}
	simpleAllTypes.Boolean = true
	simpleAllTypes.Byte = 1
	simpleAllTypes.Int = 2
	simpleAllTypes.Int16 = 3
	simpleAllTypes.Int32 = 4
	simpleAllTypes.Int64 = 5
	simpleAllTypes.Int8 = 6
	simpleAllTypes.Text = "ok"
	simpleAllTypes.Uint = 7
	simpleAllTypes.Uint16 = 8
	simpleAllTypes.Uint32 = 9
	simpleAllTypes.Uint64 = 10
	simpleAllTypes.Uint8 = 11
	bytes, err := json.Marshal(&simpleAllTypes)
	data := make(map[string]any)
	err = json.Unmarshal(bytes, &data)

	unmarshalled := new(SimpleAllTypes)
	err = Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestSimpleAllArrayTypes(t *testing.T) {
	simpleAllTypes := SimpleAllArrayTypes{}
	simpleAllTypes.Boolean = []bool{true}

	simpleAllTypes.Int = []int{2}
	simpleAllTypes.Int16 = []int16{3}
	simpleAllTypes.Int32 = []int32{4}
	simpleAllTypes.Int64 = []int64{5}

	simpleAllTypes.Text = []string{"ok"}
	simpleAllTypes.Uint = []uint{7}
	simpleAllTypes.Uint16 = []uint16{8}
	simpleAllTypes.Uint32 = []uint32{9}
	simpleAllTypes.Uint64 = []uint64{10}

	bytes, err := json.Marshal(&simpleAllTypes)
	data := make(map[string]any)
	err = json.Unmarshal(bytes, &data)

	unmarshalled := new(SimpleAllArrayTypes)
	err = Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestSimpleArray(t *testing.T) {
	data := map[string]any{
		"Number":  []any{1, 2, 3, 4, 5},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(SimpleArray)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestSimpleMap(t *testing.T) {
	data := map[string]any{
		"Number":  map[string]any{"1": 1, "2": 2},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(SimpleMap)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestSimpleStruct(t *testing.T) {
	data := map[string]any{
		"Number": map[string]any{
			"1": map[string]any{
				"Number":  10,
				"Text":    "Hello World",
				"Boolean": true,
			},
		},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(MapStruct)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestArray2D(t *testing.T) {
	data := map[string]any{
		"Number":  [][]any{{1, 2, 3}},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(Array2D)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestArrayMultiD(t *testing.T) {
	data := map[string]any{
		"Number":  [][][]any{{{1, 2, 3}}},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(ArrayMultiD)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestArrayStructMulti(t *testing.T) {
	data := map[string]any{
		"Number": [][]map[string]any{{
			{
				"Number":  10,
				"Text":    "Hello World",
				"Boolean": true,
			},
		}},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(ArrayStructMulti)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestArrayStruct(t *testing.T) {
	data := map[string]any{
		"Number": []any{
			map[string]any{

				"Number":  10,
				"Text":    "Hello World",
				"Boolean": true,
			}},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(ArrayStruct)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestArrayMap(t *testing.T) {
	data := map[string]any{
		"Number": []any{
			map[string]any{

				"Number":  10,
				"Text":    "Hello World",
				"Boolean": true,
			}},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(ArrayMap)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}

func TestArrayMap2d(t *testing.T) {
	data := map[string]any{
		"Number": [][]any{
			{
				map[string]any{

					"Number":  10,
					"Text":    "Hello World",
					"Boolean": true,
				},
			},
		},
		"Text":    "Hello World",
		"Boolean": true,
	}
	unmarshalled := new(ArrayMap2d)
	err := Unmarshal(data, unmarshalled)
	if err != nil {
		t.FailNow()
	}
}
