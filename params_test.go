package imageserver

import (
	"reflect"
	"sort"
	"testing"
)

func TestParamsSetGet(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	value, err := params.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != "bar" {
		t.Fatal("not equals")
	}
}

func TestParamsGetErrorMiss(t *testing.T) {
	params := make(Params)
	_, err := params.Get("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetString(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	value, err := params.GetString("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != "bar" {
		t.Fatal("not equals")
	}
}

func TestParamsGetStringErrorMiss(t *testing.T) {
	params := make(Params)
	_, err := params.GetString("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetStringErrorWrongType(t *testing.T) {
	params := make(Params)
	params.Set("foo", 666)
	_, err := params.GetString("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetInt(t *testing.T) {
	params := make(Params)
	params.Set("foo", 7)
	value, err := params.GetInt("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != 7 {
		t.Fatal("not equals")
	}
}

func TestParamsGetIntErrorMiss(t *testing.T) {
	params := make(Params)
	_, err := params.GetInt("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetIntErrorWrongType(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	_, err := params.GetInt("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetFloat(t *testing.T) {
	params := make(Params)
	params.Set("foo", 12.34)
	value, err := params.GetFloat("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != 12.34 {
		t.Fatal("not equals")
	}
}

func TestParamsGetFloatErrorMiss(t *testing.T) {
	params := make(Params)
	_, err := params.GetFloat("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetFloatErrorWrongType(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	_, err := params.GetFloat("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetBool(t *testing.T) {
	params := make(Params)
	params.Set("foo", true)
	value, err := params.GetBool("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != true {
		t.Fatal("Not equals")
	}
}

func TestParamsGetBoolErrorMiss(t *testing.T) {
	params := make(Params)
	_, err := params.GetBool("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetBoolErrorWrongType(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	_, err := params.GetBool("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetParams(t *testing.T) {
	params := make(Params)
	params.Set("foo", make(Params))
	_, err := params.GetParams("foo")
	if err != nil {
		t.Fatal(err)
	}
}

func TestParamsGetParamsErrorMiss(t *testing.T) {
	params := make(Params)
	_, err := params.GetParams("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsGetParamsErrorWrongType(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	_, err := params.GetParams("foo")
	testParamsCheckErrorType(err, t)
}

func TestParamsStable(t *testing.T) {
	params1 := Params{
		"a": "azerty",
		"b": []string{
			"e",
			"d",
		},
		"c": Params{
			"f": "foo",
			"g": "bar",
		},
	}

	params2 := Params{
		"c": Params{
			"g": "bar",
			"f": "foo",
		},
		"b": []string{
			"e",
			"d",
		},
		"a": "azerty",
	}

	if params1.String() != params2.String() {
		t.Fatal("not equals")
	}
}

func testParamsCheckErrorType(err error, t *testing.T) {
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*ParamError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestParamsHas(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	if !params.Has("foo") {
		t.Fatal("key does not exist")
	}
	if params.Has("xxx") {
		t.Fatal("key exists")
	}
}

func TestParamsLen(t *testing.T) {
	params := make(Params)
	params.Set("foo", "bar")
	if params.Len() != 1 {
		t.Fatal("wrong length")
	}
}

func TestParamsEmpty(t *testing.T) {
	params := make(Params)
	if !params.Empty() {
		t.Fatal("not empty")
	}
	params.Set("foo", "bar")
	if params.Empty() {
		t.Fatal("empty")
	}
}

func TestParamsKeys(t *testing.T) {
	params := make(Params)
	params.Set("b", "bar")
	params.Set("a", "foo")
	keys := params.Keys()
	sort.Strings(keys)

	expected := []string{"a", "b"}

	if !reflect.DeepEqual(keys, expected) {
		t.Fatal("not equals")
	}
}

var _ error = &ParamError{}

func TestParamError(t *testing.T) {
	err := &ParamError{Param: "param", Message: "my message"}
	_ = err.Error()
}
