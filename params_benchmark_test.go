package imageserver

import (
	"testing"
)

func BenchmarkParamsString(b *testing.B) {
	params := Params{
		"string": "foo",
		"int":    123,
		"float":  0.123,
		"params": Params{
			"foo": "bar",
			"baz": "aaaaaaaa",
		},
	}
	for i := 0; i < b.N; i++ {
		_ = params.String()
	}
}
