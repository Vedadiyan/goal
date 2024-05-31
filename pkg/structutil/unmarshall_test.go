package structutil

import "testing"

type (
	Simple struct {
		Number  int
		Text    string
		Boolean bool
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
