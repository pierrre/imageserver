package imageserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	binary.Write(buffer, binary.LittleEndian, &formatLen)
	buffer.Write(formatBytes)

	dataLen := uint32(len(image.Data))
	binary.Write(buffer, binary.LittleEndian, &dataLen)
	buffer.Write(image.Data)

	return buffer.Bytes(), nil
}

// UnmarshalBinary unserializes bytes to the Image
func (image *Image) UnmarshalBinary(data []byte) error {
	reader := bytes.NewReader(data)

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

// NewImageUnmarshalBinaryOptimized creates a new Image from serialized bytes
func NewImageUnmarshalBinaryOptimized(data []byte) (*Image, error) {
	image := new(Image)

	err := image.UnmarshalBinaryOptimized(data)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// UnmarshalBinaryOptimized unserializes bytes to the Image
func (image *Image) UnmarshalBinaryOptimized(data []byte) error {
	dataStart, dataEnd := 0, 0
	readData := func(length int) ([]byte, error) {
		dataStart = dataEnd
		dataEnd += length
		if dataEnd > len(data) {
			return nil, fmt.Errorf("unexpected end of data at index %d instead of %d", len(data), dataEnd)
		}
		return data[dataStart:dataEnd], nil
	}

	var imageFormatLength uint32
	if d, err := readData(4); err == nil {
		binary.Read(bytes.NewReader(d), binary.LittleEndian, &imageFormatLength)
	} else {
		return err
	}

	if d, err := readData(int(imageFormatLength)); err == nil {
		image.Format = string(d)
	} else {
		return err
	}

	var imageDataLength uint32
	if d, err := readData(4); err == nil {
		binary.Read(bytes.NewReader(d), binary.LittleEndian, &imageDataLength)
	} else {
		return err
	}

	if d, err := readData(int(imageDataLength)); err == nil {
		image.Data = d
	} else {
		return err
	}

	return nil
}
