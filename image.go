package imageserver

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Image represents a raw image
type Image struct {
	Format string // png, jpeg, bmp, gif, ...
	Data   []byte // raw image data
}

// NewImageUnmarshalBinary creates a new Image from serialized bytes
func NewImageUnmarshalBinary(marshalledData []byte) (*Image, error) {
	image := new(Image)

	err := image.UnmarshalBinary(marshalledData)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// MarshalBinary serializes the Image to bytes
func (image *Image) MarshalBinary() ([]byte, error) {
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

// UnmarshalBinary unserializes bytes to the Image
func (image *Image) UnmarshalBinary(marshalledData []byte) error {
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
