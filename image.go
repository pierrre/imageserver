package imageserver

import (
	"bytes"
	"encoding/gob"
)

type Image struct {
	Type string
	Data []byte
}

func (image *Image) Serialize() (serialized []byte, err error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	err = encoder.Encode(image)
	if err != nil {
		return
	}
	serialized = buffer.Bytes()
	return
}

func (image *Image) Unserialize(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(image)
	return err
}
