package imageserver

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
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

func NewImageUnmarshalBinaryExp(marshalledData []byte) (*Image, error) {
	image := new(Image)

	err := image.UnmarshalBinaryExp(marshalledData)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (image *Image) MarshalBinaryExp() ([]byte, error) {
	buffer := new(bytes.Buffer)

	formatBytes := []byte(image.Format)
	formatLen := uint32(len(formatBytes))
	err := binary.Write(buffer, binary.LittleEndian, &formatLen)
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write(formatBytes)
	if err != nil {
		return nil, err
	}

	dataLen := uint32(len(image.Data))
	err = binary.Write(buffer, binary.LittleEndian, &dataLen)
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write(image.Data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (image *Image) UnmarshalBinaryExp(marshalledData []byte) error {
	reader := bytes.NewReader(marshalledData)

	var formatLen uint32
	err := binary.Read(reader, binary.LittleEndian, &formatLen)
	if err != nil {
		return err
	}
	formatBytes := make([]byte, formatLen)
	_, err = io.ReadFull(reader, formatBytes)
	if err != nil {
		return err
	}
	image.Format = string(formatBytes)

	var dataLen uint32
	err = binary.Read(reader, binary.LittleEndian, &dataLen)
	if err != nil {
		return err
	}
	image.Data = make([]byte, dataLen)
	_, err = io.ReadFull(reader, image.Data)
	if err != nil {
		return err
	}

	return nil
}
