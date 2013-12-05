package imageserver

import (
	"reflect"
	"sort"
	"testing"
)

func TestParametersSetGet(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	value, err := parameters.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != "bar" {
		t.Fatal("not equals")
	}
}

func TestParametersHas(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	if !parameters.Has("foo") {
		t.Fatal("key does not exist")
	}
	if parameters.Has("xxx") {
		t.Fatal("key exists")
	}
}

func TestParametersLen(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	if parameters.Len() != 1 {
		t.Fatal("wrong length")
	}
}

func TestParametersEmpty(t *testing.T) {
	parameters := make(Parameters)
	if !parameters.Empty() {
		t.Fatal("not empty")
	}
	parameters.Set("foo", "bar")
	if parameters.Empty() {
		t.Fatal("empty")
	}
}

func TestParametersKeys(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("b", "bar")
	parameters.Set("a", "foo")
	keys := parameters.Keys()
	sort.Strings(keys)

	expected := []string{"a", "b"}

	if !reflect.DeepEqual(keys, expected) {
		t.Fatal("not equals")
	}
}

func TestParametersGetErrorMiss(t *testing.T) {
	parameters := make(Parameters)
	_, err := parameters.Get("foo")
	if err == nil {
		t.Fatal("no miss")
	}
}

func TestParametersGetString(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	value, err := parameters.GetString("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != "bar" {
		t.Fatal("not equals")
	}
}

func TestParametersGetStringErrorMiss(t *testing.T) {
	parameters := make(Parameters)
	_, err := parameters.GetString("foo")
	if err == nil {
		t.Fatal("no miss")
	}
}

func TestParametersGetStringErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", 666)
	_, err := parameters.GetString("foo")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParametersGetInt(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", 7)
	value, err := parameters.GetInt("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != 7 {
		t.Fatal("not equals")
	}
}

func TestParametersGetIntErrorMiss(t *testing.T) {
	parameters := make(Parameters)
	_, err := parameters.GetInt("foo")
	if err == nil {
		t.Fatal("no miss")
	}
}

func TestParametersGetIntErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	_, err := parameters.GetInt("foo")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParametersGetBool(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", true)
	value, err := parameters.GetBool("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != true {
		t.Fatal("Not equals")
	}
}

func TestParametersGetBoolErrorMiss(t *testing.T) {
	parameters := make(Parameters)
	_, err := parameters.GetBool("foo")
	if err == nil {
		t.Fatal("no miss")
	}
}

func TestParametersGetBoolErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	_, err := parameters.GetBool("foo")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParametersGetParameters(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", make(Parameters))
	_, err := parameters.GetParameters("foo")
	if err != nil {
		t.Fatal(err)
	}
}

func TestParametersGetParametersErrorMiss(t *testing.T) {
	parameters := make(Parameters)
	_, err := parameters.GetParameters("foo")
	if err == nil {
		t.Fatal("no miss")
	}
}

func TestParametersGetParametersErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	_, err := parameters.GetParameters("foo")
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParametersStable(t *testing.T) {
	parameters1 := Parameters{
		"a": "azerty",
		"b": []string{
			"e",
			"d",
		},
		"c": Parameters{
			"f": "foo",
			"g": "bar",
		},
	}

	parameters2 := Parameters{
		"c": Parameters{
			"g": "bar",
			"f": "foo",
		},
		"b": []string{
			"e",
			"d",
		},
		"a": "azerty",
	}

	if parameters1.String() != parameters2.String() {
		t.Fatal("not equals")
	}
}
