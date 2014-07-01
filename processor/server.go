package processor

import (
	"github.com/pierrre/imageserver"
)

// ProcessorImageServer is a Processor ImageServer
type ProcessorImageServer struct {
	ImageServer imageserver.ImageServer
	Processor   Processor
}

// Get gets an Image from the underlying ImageServer, then processes it with the Processor
func (pis *ProcessorImageServer) Get(parameters imageserver.Parameters) (*imageserver.Image, error) {
	image, err := pis.ImageServer.Get(parameters)
	if err != nil {
		return nil, err
	}

	image, err = pis.Processor.Process(image, parameters)
	if err != nil {
		return nil, err
	}

	return image, err
}
