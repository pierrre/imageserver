package imageserver_test

import (
	"testing"

	. "github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkImageMarshalBinary(b *testing.B) {
	for _, tc := range []struct {
		name string
		im   *Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
		{"Huge", testdata.Huge},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := tc.im.MarshalBinary()
				if err != nil {
					b.Fatal(err)
				}
			}
			b.SetBytes(int64(len(tc.im.Data)))
		})
	}
}

func BenchmarkImageUnmarshalBinary(b *testing.B) {
	for _, tcnc := range []struct {
		name   string
		nocopy bool
	}{
		{"Normal", false},
		{"NoCopy", true},
	} {
		b.Run(tcnc.name, func(b *testing.B) {
			for _, tcim := range []struct {
				name string
				im   *Image
			}{
				{"Small", testdata.Small},
				{"Medium", testdata.Medium},
				{"Large", testdata.Large},
				{"Huge", testdata.Huge},
			} {
				b.Run(tcim.name, func(b *testing.B) {
					data, err := tcim.im.MarshalBinary()
					if err != nil {
						b.Fatal(err)
					}
					imNew := new(Image)
					var m func([]byte) error
					if tcnc.nocopy {
						m = imNew.UnmarshalBinaryNoCopy
					} else {
						m = imNew.UnmarshalBinary
					}
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						err := m(data)
						if err != nil {
							b.Fatal(err)
						}
					}
					b.SetBytes(int64(len(tcim.im.Data)))
				})
			}
		})
	}
}
