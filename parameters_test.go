package imageserver

import (
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

func TestParametersGetParametersErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	_, err := parameters.GetParameters("foo")
	if err == nil {
		t.Fatal("no error")
	}
}
