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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = params.String()
		}
	})
}
