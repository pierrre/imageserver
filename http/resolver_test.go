package http

import (
	"testing"
)

func TestResolverFuncInterface(t *testing.T) {
	var _ Resolver = ResolverFunc(nil)
}

func TestResolverFunc(t *testing.T) {
	rf := ResolverFunc(func(parameter string) string {
		return ""
	})
	rf.Resolve("")
}

func TestResolve(t *testing.T) {
	rf := ResolverFunc(func(parameter string) string {
		return "foobar"
	})
	httpParameter := Resolve(rf, "test")
	if httpParameter != "foobar" {
		t.Fatal("not equals")
	}
}

func TestResolveNotResolver(t *testing.T) {
	httpParameter := Resolve(nil, "test")
	if httpParameter != "" {
		t.Fatal("not equals")
	}
}
