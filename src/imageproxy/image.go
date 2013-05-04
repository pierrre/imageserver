package imageproxy

import (
	"bytes"
	"encoding/gob"
)

type Image struct {
	Type string
	Data []byte
}

func (image *Image) serialize() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(image)
	return buffer.Bytes()
}

func (image *Image) unserialize(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(image)
}
