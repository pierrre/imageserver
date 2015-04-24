package imageserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

/*
Image represents a raw image.


Binary encoding:

- Format length (uint32)

- Format (string)

- Data length (uint32)

- Data([]byte)

Numbers are encoded using little-endian order.
*/
type Image struct {
	Format string // png, jpeg, bmp, gif, ...
	Data   []byte // raw image data
}

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// MarshalBinary implements encoding.BinaryMarshaler.
//
// It's very unlikely that it returns an error. (impossible?)
func (im *Image) MarshalBinary() ([]byte, error) {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	formatLen := uint32(len(im.Format))
	binary.Write(buf, binary.LittleEndian, &formatLen)
	buf.Write([]byte(im.Format))

	dataLen := uint32(len(im.Data))
	binary.Write(buf, binary.LittleEndian, &dataLen)
	buf.Write(im.Data)

	b := buf.Bytes()
	bufferPool.Put(buf)
	return b, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
//
// It copies data then call UnmarshalBinaryNoCopy().
func (im *Image) UnmarshalBinary(data []byte) error {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	return im.UnmarshalBinaryNoCopy(dataCopy)
}

// UnmarshalBinaryNoCopy is like encoding.BinaryUnmarshaler but does no copy.
//
// The caller must not reuse data after that.
func (im *Image) UnmarshalBinaryNoCopy(data []byte) error {
	dataPosition := 0
	readData := func(length int) ([]byte, error) {
		dataEnd := dataPosition + length
		if dataEnd > len(data) {
			return nil, &ImageError{Message: fmt.Sprintf("unmarshal: unexpected end of data at index %d instead of %d", len(data), dataEnd)}
		}
		d := data[dataPosition:dataEnd]
		dataPosition = dataEnd
		return d, nil
	}

	var formatLen uint32
	if d, err := readData(4); err == nil {
		binary.Read(bytes.NewReader(d), binary.LittleEndian, &formatLen)
	} else {
		return err
	}

	if d, err := readData(int(formatLen)); err == nil {
		im.Format = string(d)
	} else {
		return err
	}

	var dataLen uint32
	if d, err := readData(4); err == nil {
		binary.Read(bytes.NewReader(d), binary.LittleEndian, &dataLen)
	} else {
		return err
	}

	if d, err := readData(int(dataLen)); err == nil {
		im.Data = d
	} else {
		return err
	}

	return nil
}

// ImageEqual compares two images and returns true if they are equal.
func ImageEqual(im1, im2 *Image) bool {
	if im1 == im2 {
		return true
	}
	if im1 == nil || im2 == nil {
		return false
	}
	if im1.Format != im2.Format {
		return false
	}
	if !bytes.Equal(im1.Data, im2.Data) {
		return false
	}
	return true
}

// ImageError is an Image error.
type ImageError struct {
	Message string
}

func (err *ImageError) Error() string {
	return fmt.Sprintf("image error: %s", err.Message)
}
