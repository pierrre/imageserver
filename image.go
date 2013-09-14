package imageserver

import (
	"bytes"
	"encoding/gob"
)

// Internal image structure
//
// An image is composed of a type (png, jpeg, bmp, gif, ...) and data (slice of byte)
//
// This data structure is easy to serialize (cache) and manipulate (processor implementation)
type Image struct {
	Type string
	Data []byte
}

// Serialize image to bytes
func (image *Image) Marshal() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(image); err != nil {
		return nil, err
	}
	data := buffer.Bytes()
	return data, nil
}

// Fill image with serialized bytes
func (image *Image) Unmarshal(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(image)
	return err
}
