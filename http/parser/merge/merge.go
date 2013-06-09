package merge

import (
	"github.com/pierrre/imageserver"
	imageserver_http "github.com/pierrre/imageserver/http"
	"net/http"
)

type MergeParser []imageserver_http.Parser

func (parser MergeParser) Parse(request *http.Request) (parameters imageserver.Parameters, err error) {
	for _, subParser := range parser {
		parameters, err = parseAndMerge(request, subParser, parameters)
		if err != nil {
			return
		}

	}
	return
}

func parseAndMerge(request *http.Request, parser imageserver_http.Parser, inParameters imageserver.Parameters) (parameters imageserver.Parameters, err error) {
	parameters = inParameters
	subParameters, err := parser.Parse(request)
	if err != nil {
		return
	}
	if parameters != nil {
		for k, v := range subParameters {
			parameters[k] = v
		}
	} else {
		parameters = subParameters
	}
	return
}
