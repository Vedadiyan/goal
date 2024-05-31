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
	SimpleMap struct {
		Number  map[string]int
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
