package file

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_cache "github.com/pierrre/imageserver/cache"
	cachetest "github.com/pierrre/imageserver/cache/_test"
	"github.com/pierrre/imageserver/testdata"
	"io/ioutil"
	"os"
)

var _ imageserver_cache.Cache = &Cache{}
var testDirPath string

func TestMain(m *testing.M) {
	os.Exit(realMain(m))
}

func realMain(m *testing.M) int {
	path, err := ioutil.TempDir("", "filecache")
	if err != nil {
		return 1
	}
	testDirPath = path
	defer func() {
		_ = os.RemoveAll(testDirPath)
	}()
	return m.Run()
}

func TestGetSet(t *testing.T) {
	cache := newTestCache(t)
	cachetest.TestGetSet(t, cache)
}

func TestGetMiss(t *testing.T) {
	cache := newTestCache(t)
	cachetest.TestGetMiss(t, cache)
}

func TestPathIsNotSet(t *testing.T) {
	cache := &Cache{Path: ""}
	_, err := cache.Get(cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestFileExistsButCantRead(t *testing.T) {
	cache := newTestCache(t)
	if err := os.Chmod(filepath.Join(testDirPath, cachetest.KeyValid), 0111); err != nil {
		t.Fatal("os.Chmod return error.")
	}
	defer func() {
		err := os.Chmod(filepath.Join(testDirPath, cachetest.KeyValid), 0644)
		if err != nil {
			t.Fatal("os.Chmod return error.")
		}
	}()
	_, err := cache.Get(cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
}

func TestGetErrorUnmarshal(t *testing.T) {
	cache := newTestCache(t)
	data, err := testdata.Medium.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	data = data[:len(data)-1]
	err = cache.setData(cachetest.KeyValid, data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cache.Get(cachetest.KeyValid, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func TestSetErrorMarshal(t *testing.T) {
	cache := newTestCache(t)
	im := &imageserver.Image{
		Format: strings.Repeat("a", imageserver.ImageFormatMaxLen+1),
	}
	err := cache.Set(cachetest.KeyValid, im, imageserver.Params{})
	if err == nil {
		t.Fatal("no error")
	}
	if _, ok := err.(*imageserver.ImageError); !ok {
		t.Fatalf("unexpected error type: %T", err)
	}
}

func newTestCache(tb testing.TB) *Cache {
	cache := &Cache{Path: testDirPath}
	return cache
}
