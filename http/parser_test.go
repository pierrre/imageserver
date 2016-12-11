package http

import (
	"net/http"
	"testing"

	"github.com/pierrre/imageserver"
	imageserver_source "github.com/pierrre/imageserver/source"
)

var _ Parser = ListParser{}

func TestListParserParse(t *testing.T) {
	parser := ListParser{
		&SourceParser{},
	}
	req, err := http.NewRequest("GET", "http://localhost?source=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if !params.Has(imageserver_source.Param) {
		t.Fatal("not set")
	}
}

func TestListParserParseError(t *testing.T) {
	parser := ListParser{
		&testErrorParser{},
	}
	req, err := http.NewRequest("GET", "http://localhost?error=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestListParserResolve(t *testing.T) {
	parser := ListParser{
		&SourceParser{},
	}

	httpParam := parser.Resolve(imageserver_source.Param)
	if httpParam != imageserver_source.Param {
		t.Fatal("not equals")
	}

	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourceParser{}

func TestSourceParserParse(t *testing.T) {
	parser := &SourceParser{}
	req, err := http.NewRequest("GET", "http://localhost?source=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	source, err := params.GetString(imageserver_source.Param)
	if err != nil {
		t.Fatal(err)
	}
	if source != "foo" {
		t.Fatal("not equals")
	}
}

func TestSourceParserParseUndefined(t *testing.T) {
	parser := &SourceParser{}
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}

	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has(imageserver_source.Param) {
		t.Fatal("should not be set")
	}
}

func TestSourceParserResolve(t *testing.T) {
	parser := &SourceParser{}
	httpParam := parser.Resolve(imageserver_source.Param)
	if httpParam != imageserver_source.Param {
		t.Fatal("not equals")
	}
	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourcePathParser{}

func TestSourcePathParserParse(t *testing.T) {
	parser := &SourcePathParser{}
	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = parser.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	source, err := params.GetString(imageserver_source.Param)
	if err != nil {
		t.Fatal(err)
	}
	if source != "/foobar" {
		t.Fatal("not equals")
	}
}

func TestSourcePathParserResolve(t *testing.T) {
	parser := &SourcePathParser{}
	httpParam := parser.Resolve(imageserver_source.Param)
	if httpParam != "path" {
		t.Fatal("not equals")
	}
	httpParam = parser.Resolve("foobar")
	if httpParam != "" {
		t.Fatal("not equals")
	}
}

var _ Parser = &SourceTransformParser{}

func TestSourceTransformParser(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?source=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	ps := &SourceTransformParser{
		Parser: &SourceParser{},
		Transform: func(source string) string {
			return "bar"
		},
	}
	params := imageserver.Params{}
	err = ps.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	source, err := params.GetString(imageserver_source.Param)
	if err != nil {
		t.Fatal(err)
	}
	if source != "bar" {
		t.Fatal("not equals")
	}
}

func TestSourceTransformParserUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	ps := &SourceTransformParser{
		Parser: &SourceParser{},
	}
	params := imageserver.Params{}
	err = ps.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has(imageserver_source.Param) {
		t.Fatal("should not be set")
	}
}

func TestSourceTransformParserErrorParse(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?error=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	ps := &SourceTransformParser{
		Parser: &testErrorParser{},
	}
	params := imageserver.Params{}
	err = ps.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestSourceTransformParserErrorParams(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	ps := &SourceTransformParser{
		Parser: &SourceParser{},
	}
	params := imageserver.Params{
		imageserver_source.Param: 666,
	}
	err = ps.Parse(req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

var _ Parser = &SourcePrefixParser{}

func TestSourcePrefixParser(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?source=bar", nil)
	if err != nil {
		t.Fatal(err)
	}
	ps := &SourcePrefixParser{
		Parser: &SourceParser{},
		Prefix: "foo",
	}
	params := imageserver.Params{}
	err = ps.Parse(req, params)
	if err != nil {
		t.Fatal(err)
	}
	source, err := params.GetString(imageserver_source.Param)
	if err != nil {
		t.Fatal(err)
	}
	if source != "foobar" {
		t.Fatal("not equals")
	}
}

func TestParseQueryString(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?string=foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	ParseQueryString("string", req, params)
	s, err := params.GetString("string")
	if err != nil {
		t.Fatal(err)
	}
	if s != "foo" {
		t.Fatal("not equals")
	}
}

func TestParseQueryStringUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	ParseQueryString("string", req, params)
	if params.Has("string") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryInt(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?int=42", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt("int", req, params)
	if err != nil {
		t.Fatal(err)
	}
	i, err := params.GetInt("int")
	if err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal("not equals")
	}
}

func TestParseQueryIntUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt("int", req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("int") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryIntError(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?int=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt("int", req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseQueryInt64(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?int64=42", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt64("int64", req, params)
	if err != nil {
		t.Fatal(err)
	}
	i, err := params.GetInt64("int64")
	if err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Fatal("not equals")
	}
}

func TestParseQueryInt64Undefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt64("int64", req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("int64") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryInt64Error(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?int64=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryInt64("int64", req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseQueryFloat(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?float=12.34", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryFloat("float", req, params)
	if err != nil {
		t.Fatal(err)
	}
	f, err := params.GetFloat("float")
	if err != nil {
		t.Fatal(err)
	}
	if f != 12.34 {
		t.Fatal("not equals")
	}
}

func TestParseQueryFloatUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryFloat("float", req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("float") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryFloatError(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?float=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryFloat("float", req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

func TestParseQueryBool(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?bool=true", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryBool("bool", req, params)
	if err != nil {
		t.Fatal(err)
	}
	b, err := params.GetBool("bool")
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Fatal("not equals")
	}
}

func TestParseQueryBoolUndefined(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryBool("bool", req, params)
	if err != nil {
		t.Fatal(err)
	}
	if params.Has("bool") {
		t.Fatal("should not be set")
	}
}

func TestParseQueryBoolError(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost?bool=invalid", nil)
	if err != nil {
		t.Fatal(err)
	}
	params := imageserver.Params{}
	err = ParseQueryBool("bool", req, params)
	if err == nil {
		t.Fatal("no error")
	}
}

type testErrorParser struct{}

func (prs *testErrorParser) Parse(req *http.Request, params imageserver.Params) error {
	v := req.URL.Query().Get("error")
	if v != "" {
		return &imageserver.ParamError{
			Param:   "error",
			Message: v,
		}
	}
	return nil
}

func (prs *testErrorParser) Resolve(param string) string {
	if param == "error" {
		return "error"
	}
	return ""
}
