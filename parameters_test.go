package imageserver

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	value, err := parameters.Get("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != "bar" {
		t.Fatal("Not equals")
	}
}

func TestGetErrorNotFound(t *testing.T) {
	parameters := make(Parameters)
	_, err := parameters.Get("foo")
	if err == nil {
		t.Fatal(err)
	}
}

func TestGetString(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	value, err := parameters.GetString("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != "bar" {
		t.Fatal("Not equals")
	}
}

func TestGetStringErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", 666)
	_, err := parameters.GetString("foo")
	if err == nil {
		t.Fatal("No error")
	}
}

func TestGetInt(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", 7)
	value, err := parameters.GetInt("foo")
	if err != nil {
		t.Fatal(err)
	}
	if value != 7 {
		t.Fatal("Not equals")
	}
}

func TestGetIntErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	_, err := parameters.GetInt("foo")
	if err == nil {
		t.Fatal("No error")
	}
}

func TestGetBool(t *testing.T) {
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

func TestGetBoolErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	_, err := parameters.GetBool("foo")
	if err == nil {
		t.Fatal("No error")
	}
}

func TestGetParameters(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", make(Parameters))
	_, err := parameters.GetParameters("foo")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetParametersErrorWrongType(t *testing.T) {
	parameters := make(Parameters)
	parameters.Set("foo", "bar")
	_, err := parameters.GetParameters("foo")
	if err == nil {
		t.Fatal("No error")
	}
}
