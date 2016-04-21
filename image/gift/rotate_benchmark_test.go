package gift

import (
	"image/color"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_image "github.com/pierrre/imageserver/image"
	"github.com/pierrre/imageserver/testdata"
)

func BenchmarkRotateProcessor90(b *testing.B) {
	benchmarkRotateProcessorRotation(b, 90)
}

func BenchmarkRotateProcessor180(b *testing.B) {
	benchmarkRotateProcessorRotation(b, 180)
}

func BenchmarkRotateProcessor270(b *testing.B) {
	benchmarkRotateProcessorRotation(b, 270)
}

func BenchmarkRotateProcessor45(b *testing.B) {
	benchmarkRotateProcessorRotation(b, 45)
}

func benchmarkRotateProcessorRotation(b *testing.B, rot float64) {
	benchmarkRotateProcessor(b, imageserver.Params{"rotation": rot})
}

func BenchmarkRotateProcessorInterpolationNearestNeighbor(b *testing.B) {
	benchmarkRotateProcessorInterpolation(b, "nearest_neighbor")
}

func BenchmarkRotateProcessorInterpolationLinear(b *testing.B) {
	benchmarkRotateProcessorInterpolation(b, "linear")
}

func BenchmarkRotateProcessorInterpolationCubic(b *testing.B) {
	benchmarkRotateProcessorInterpolation(b, "cubic")
}

func benchmarkRotateProcessorInterpolation(b *testing.B, interp string) {
	benchmarkRotateProcessor(b, imageserver.Params{
		"rotation":      45.0,
		"interpolation": interp,
	})
}

func benchmarkRotateProcessor(b *testing.B, params imageserver.Params) {
	nim, err := imageserver_image.Decode(testdata.Medium)
	if err != nil {
		b.Fatal(err)
	}
	params = imageserver.Params{
		rotateParam: params,
	}
	prc := &RotateProcessor{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := prc.Process(nim, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseHexColor3(b *testing.B) {
	benchmarkParseHexColor(b, "FFF")
}

func BenchmarkParseHexColor4(b *testing.B) {
	benchmarkParseHexColor(b, "FFFF")
}

func BenchmarkParseHexColor6(b *testing.B) {
	benchmarkParseHexColor(b, "FFFFFF")
}

func BenchmarkParseHexColor8(b *testing.B) {
	benchmarkParseHexColor(b, "FFFFFFFF")
}

var benchmarkParseHexColorRes color.Color

func benchmarkParseHexColor(b *testing.B, s string) {
	var res color.Color
	var err error
	for i := 0; i < b.N; i++ {
		res, err = parseHexColor(s)
		if err != nil {
			b.Fatal(err)
		}
	}
	benchmarkParseHexColorRes = res
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
