package imageserver

import (
	"bytes"
	"encoding/gob"
)

// Image represents a raw image
type Image struct {
	Format string // png, jpeg, bmp, gif, ...
	Data   []byte // raw image data
}

// NewImageUnmarshal creates a new Image from serialized bytes
func NewImageUnmarshal(marshalledData []byte) (*Image, error) {
	image := new(Image)

	err := image.Unmarshal(marshalledData)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// Marshal serializes the Image to bytes
func (image *Image) Marshal() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(image); err != nil {
		return nil, err
	}
	data := buffer.Bytes()
	return data, nil
}

// Unmarshal unserializes bytes to the Image
func (image *Image) Unmarshal(marshalledData []byte) error {
	buffer := bytes.NewBuffer(marshalledData)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(image)
	return err
}
