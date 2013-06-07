package merge

import (
	"github.com/pierrre/imageproxy"
	"net/http"
)

type MergeRequestParser []imageproxy.RequestParser

func (parser MergeRequestParser) ParseRequest(request *http.Request) (parameters imageproxy.Parameters, err error) {
	for _, subParser := range parser {
		parameters, err = parseAndMerge(request, subParser, parameters)
		if err != nil {
			return
		}

	}
	return
}

func parseAndMerge(request *http.Request, parser imageproxy.RequestParser, inParameters imageproxy.Parameters) (parameters imageproxy.Parameters, err error) {
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
