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
func (image *Image) Marshal() (data []byte, err error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	if err = encoder.Encode(image); err != nil {
		return
	}
	data = buffer.Bytes()
	return
}

// Fill image with serialized bytes
func (image *Image) Unmarshal(data []byte) (err error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(image)
	return
}
