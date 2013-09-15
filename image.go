package imageserver

import (
	"bytes"
	"encoding/gob"
)

type Image struct {
	Type string // png, jpeg, bmp, gif, ...
	Data []byte // raw image data
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
