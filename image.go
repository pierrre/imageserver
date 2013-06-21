package imageserver

import (
	"bytes"
	"encoding/gob"
)

type Image struct {
	Type string
	Data []byte
}

func (image *Image) Marshal() (data []byte, err error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	if err = encoder.Encode(image); err != nil {
		return
	}
	data = buffer.Bytes()
	return
}

func (image *Image) Unmarshal(data []byte) (err error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(image)
	return
}
