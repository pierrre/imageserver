package jpeg

import (
	"strconv"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image_test "github.com/pierrre/imageserver/image/_test"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkSize(b *testing.B) {
	params := imageserver.Params{}
	for _, tc := range []struct {
		name string
		im   *imageserver.Image
	}{
		{"Small", testdata.Small},
		{"Medium", testdata.Medium},
		{"Large", testdata.Large},
		{"Huge", testdata.Huge},
	} {
		benchmark(b, tc.name, tc.im, params)
	}
}

func BenchmarkQuality(b *testing.B) {
	for _, q := range []int{
		1, 25, 50, 75, 85, 90, 95, 100,
	} {
		benchmark(b, strconv.Itoa(q), testdata.Medium, imageserver.Params{
			"quality": q,
		})
	}
}

func benchmark(b *testing.B, name string, im *imageserver.Image, params imageserver.Params) {
	b.Run(name, func(b *testing.B) {
		imageserver_image_test.BenchmarkEncoder(b, &Encoder{}, im, params)
	})
}
