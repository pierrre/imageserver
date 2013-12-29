package imageserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
//
// It's very unlikely that it returns an error (impossible?)
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

// ImageEqual compares two images and returns true if they are equal
func ImageEqual(image1, image2 *Image) bool {
	if image1 == image2 {
		return true
	}

	if image1 == nil || image2 == nil {
		return false
	}

	if image1.Format != image2.Format {
		return false
	}

	if !bytes.Equal(image1.Data, image2.Data) {
		return false
	}

	return true
}
