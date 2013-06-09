package merge

import (
	"github.com/pierrre/imageserver"
	"net/http"
)

type MergeRequestParser []imageserver.RequestParser

func (parser MergeRequestParser) ParseRequest(request *http.Request) (parameters imageserver.Parameters, err error) {
	for _, subParser := range parser {
		parameters, err = parseAndMerge(request, subParser, parameters)
		if err != nil {
			return
		}

	}
	return
}

func parseAndMerge(request *http.Request, parser imageserver.RequestParser, inParameters imageserver.Parameters) (parameters imageserver.Parameters, err error) {
	parameters = inParameters
	subParameters, err := parser.ParseRequest(request)
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
