package gift

import (
	"context"
	"fmt"
	"image/color"
	"strconv"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkRotateProcessorRotation(b *testing.B) {
	for _, r := range []float64{
		90, 180, 270, 45,
	} {
		benchmarkRotateProcessor(b, fmt.Sprint(r), imageserver.Params{"rotation": r})
	}
}

func BenchmarkRotateProcessorInterpolation(b *testing.B) {
	for _, it := range []string{
		"nearest_neighbor",
		"linear",
		"cubic",
	} {
		benchmarkRotateProcessor(b, it, imageserver.Params{
			"rotation":      45.0,
			"interpolation": it,
		})
	}
}

func benchmarkRotateProcessor(b *testing.B, name string, params imageserver.Params) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	params = imageserver.Params{
		rotateParam: params,
	}
	prc := &RotateProcessor{}
	ctx := context.Background()
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := prc.Process(ctx, nim, params)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

var benchmarkParseHexColorRes color.Color

func BenchmarkParseHexColor(b *testing.B) {
	for _, s := range []string{
		"FFF",
		"FFFF",
		"FFFFFF",
		"FFFFFFFF",
	} {
		b.Run(strconv.Itoa(len(s)), func(b *testing.B) {
			var res color.Color
			var err error
			for i := 0; i < b.N; i++ {
				res, err = parseHexColor(s)
				if err != nil {
					b.Fatal(err)
				}
			}
			benchmarkParseHexColorRes = res
		})
	}
}

var benchmarkHexStringToIntsRes []uint8

func BenchmarkHexStringToInts(b *testing.B) {
	s := "FFFFFFFF"
	var res []uint8
	var err error
	for i := 0; i < b.N; i++ {
		res, err = hexStringToInts(s)
		if err != nil {
			b.Fatal(err)
		}
	}
	benchmarkHexStringToIntsRes = res
}
