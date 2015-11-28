package graphicsmagick

import (
	"container/list"
	"strconv"
	"testing"
	"time"

	"github.com/pierrre/imageserver"
	"github.com/pierrre/imageserver/testdata"
)

var _ imageserver.Server = &Server{}

func TestGet(t *testing.T) {
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := server.Get(params)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetErrorTimeout(t *testing.T) {
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
		Timeout:    1 * time.Nanosecond,
	}
	params := imageserver.Params{
		globalParam: imageserver.Params{
			"width":  100,
			"height": 100,
		},
	}
	_, err := server.Get(params)
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

type testSample struct {
	Params            imageserver.Params
	ExpectedArguments string
}

func convertArgumentsToString(list *list.List) string {
	var str string
	for e := list.Front(); e != nil; e = e.Next() {
		substr, ok := e.Value.(string)
		if ok {
			str += substr + " "
		} else {
			return ""
		}
	}
	return str
}

func TestBuildArgumentsResize(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"width":  100,
				"height": 150,
			},
			ExpectedArguments: "-resize 100x150 ",
		},

		testSample{
			Params: imageserver.Params{
				"width":  200,
				"height": 120,
				"fill":   true,
			},
			ExpectedArguments: "-resize 200x120^ ",
		},

		testSample{
			Params: imageserver.Params{
				"width":        120,
				"height":       100,
				"ignore_ratio": true,
			},
			ExpectedArguments: "-resize 120x100! ",
		},

		testSample{
			Params: imageserver.Params{
				"width":              120,
				"height":             100,
				"only_shrink_larger": true,
			},
			ExpectedArguments: "-resize 120x100> ",
		},

		testSample{
			Params: imageserver.Params{
				"width":                120,
				"height":               100,
				"only_enlarge_smaller": true,
			},
			ExpectedArguments: "-resize 120x100< ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		_, _, err := server.buildArgumentsResize(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}

}

func TestBuildArgumentsBackground(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"background": "59ecaf",
			},
			ExpectedArguments: "-background #59ecaf ",
		},

		testSample{
			Params: imageserver.Params{
				"background": "0059ecaf",
			},
			ExpectedArguments: "-background #0059ecaf ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsBackground(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}

}
func TestBuildArgumentsExtent(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"width":  100,
				"height": 150,
				"extent": true,
			},
			ExpectedArguments: "-resize 100x150 -extent 100x150 ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		width, height, err := server.buildArgumentsResize(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		err = server.buildArgumentsExtent(arguments, ts.Params, width, height)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}

}

func TestBuildArgumentsQuality(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"quality": 70,
				"format":  "jpeg",
			},
			ExpectedArguments: "-quality 70 ",
		},
		testSample{
			Params: imageserver.Params{
				"quality": 120,
				"format":  "png",
			},
			ExpectedArguments: "-quality 120 ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		format, err := ts.Params.GetString("format")
		if err != nil {
			t.Errorf(err.Error())
		}
		err = server.buildArgumentsQuality(arguments, ts.Params, format)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}

}

func TestBuildArgumentsGravity(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"gravity": "n",
			},
			ExpectedArguments: "-gravity North ",
		},
		testSample{
			Params: imageserver.Params{
				"gravity": "s",
			},
			ExpectedArguments: "-gravity South ",
		},
		testSample{
			Params: imageserver.Params{
				"gravity": "e",
			},
			ExpectedArguments: "-gravity East ",
		},
		testSample{
			Params: imageserver.Params{
				"gravity": "w",
			},
			ExpectedArguments: "-gravity West ",
		},
		testSample{
			Params: imageserver.Params{
				"gravity": "ne",
			},
			ExpectedArguments: "-gravity NorthEast ",
		},
		testSample{
			Params: imageserver.Params{
				"gravity": "nw",
			},
			ExpectedArguments: "-gravity NorthWest ",
		},
		testSample{
			Params: imageserver.Params{
				"gravity": "se",
			},
			ExpectedArguments: "-gravity SouthEast ",
		},
		testSample{
			Params: imageserver.Params{
				"gravity": "sw",
			},
			ExpectedArguments: "-gravity SouthWest ",
		},
		testSample{
			Params:            imageserver.Params{},
			ExpectedArguments: "-gravity Center ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsGravity(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}

}

func TestBuildArgumentsCrop(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"crop": "100,150",
			},
			ExpectedArguments: "-crop 100x150 +repage ",
		},
		testSample{
			Params: imageserver.Params{
				"crop": "120,130,40,50",
			},
			ExpectedArguments: "-crop 120x130+40+50 +repage ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsCrop(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}
}

func TestBuildArgumentsRotate(t *testing.T) {
	ts := testSample{
		Params: imageserver.Params{
			"rotate": 90,
		},
		ExpectedArguments: "-rotate 90 ",
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	arguments := list.New()
	err := server.buildArgumentsRotate(arguments, ts.Params)
	if err != nil {
		t.Errorf(err.Error())
	}
	args := convertArgumentsToString(arguments)
	if args != ts.ExpectedArguments {
		t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
	}
}

func TestBuildArgumentsMonochrome(t *testing.T) {
	ts := testSample{
		Params: imageserver.Params{
			"monochrome": true,
		},
		ExpectedArguments: "-monochrome ",
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	arguments := list.New()
	err := server.buildArgumentsMonochrome(arguments, ts.Params)
	if err != nil {
		t.Errorf(err.Error())
	}
	args := convertArgumentsToString(arguments)
	if args != ts.ExpectedArguments {
		t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
	}
}

func TestBuildArgumentsGrey(t *testing.T) {
	ts := testSample{
		Params: imageserver.Params{
			"grey": true,
		},
		ExpectedArguments: "-colorspace GRAY ",
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	arguments := list.New()
	err := server.buildArgumentsGrey(arguments, ts.Params)
	if err != nil {
		t.Errorf(err.Error())
	}
	args := convertArgumentsToString(arguments)
	if args != ts.ExpectedArguments {
		t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
	}
}

func TestBuildArgumentsStrip(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"no_strip": true,
			},
			ExpectedArguments: "",
		},
		testSample{
			Params: imageserver.Params{
				"no_strip": false,
			},
			ExpectedArguments: "-strip ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsStrip(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}
}

func TestBuildArgumentsTrim(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"trim": true,
			},
			ExpectedArguments: "-trim ",
		},
		testSample{
			Params: imageserver.Params{
				"trim": false,
			},
			ExpectedArguments: "",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsTrim(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}
}

func TestBuildArgumentsInterlace(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"no_interlace": true,
			},
			ExpectedArguments: "",
		},
		testSample{
			Params: imageserver.Params{
				"no_interlace": false,
			},
			ExpectedArguments: "-interlace Line ",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsInterlace(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}
}

func TestBuildArgumentsFlip(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"flip": true,
			},
			ExpectedArguments: "-flip ",
		},
		testSample{
			Params: imageserver.Params{
				"flip": false,
			},
			ExpectedArguments: "",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsFlip(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}
}

func TestBuildArgumentsFlop(t *testing.T) {
	testSamples := []testSample{
		testSample{
			Params: imageserver.Params{
				"flop": true,
			},
			ExpectedArguments: "-flop ",
		},
		testSample{
			Params: imageserver.Params{
				"flop": false,
			},
			ExpectedArguments: "",
		},
	}
	server := &Server{
		Server:     &imageserver.StaticServer{Image: testdata.Medium},
		Executable: "gm",
	}
	for _, ts := range testSamples {
		arguments := list.New()
		err := server.buildArgumentsFlop(arguments, ts.Params)
		if err != nil {
			t.Errorf(err.Error())
		}
		args := convertArgumentsToString(arguments)
		if args != ts.ExpectedArguments {
			t.Errorf("Failed to build arguments for gm.\nGot:  " + args + "\nWant: " + ts.ExpectedArguments)
		}
	}
}

func TestConvertArgumentsToString(t *testing.T) {
	arguments := list.New()
	arguments.PushBack("a0")
	arguments.PushBack("a1")
	arguments.PushBack("a2")
	arguments.PushBack("a3")

	argumentsString := convertArgumentsToString(arguments)

	if argumentsString != "a0 a1 a2 a3 " {
		t.Errorf("Failed to convert arguments to string")
	}
}

func TestConvertArgumentsToSlice(t *testing.T) {
	arguments := list.New()
	arguments.PushBack("a0")
	arguments.PushBack("a1")
	arguments.PushBack("a2")
	arguments.PushBack("a3")

	argumentSlice := convertArgumentsToSlice(arguments)

	for i, e := range argumentSlice {
		if e != "a"+strconv.Itoa(i) {
			t.Errorf("Failed to convert arguments to slice")
			break
		}
	}
}
