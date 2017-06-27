package imageserver

import (
	"encoding/binary"
	"fmt"
)

const (
	// ImageFormatMaxLen is the maximum length for the Image's format.
	ImageFormatMaxLen = 1 << 8 // 256 B
	// ImageDataMaxLen is the maximum length for the Image's data.
	ImageDataMaxLen = 1 << 30 // 1 GiB
)

var (
	imageByteOrder = binary.LittleEndian
)

// Image is a raw image.
//
// Binary encoding:
//  - Format length (uint32)
//  - Format (string)
//  - Data length (uint32)
//  - Data([]byte)
// Numbers are encoded using little-endian order.
type Image struct {
	// Format is the format used to encode the image.
	//
	// e.g. png, jpeg, bmp, gif, ...
	Format string

	// Data contains the raw data of the encoded image.
	Data []byte
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (im *Image) MarshalBinary() ([]byte, error) {
	if len(im.Format) > ImageFormatMaxLen {
		return nil, &ImageError{Message: fmt.Sprintf("marshal: format length %d is greater than the maximum value %d", len(im.Format), ImageFormatMaxLen)}
	}
	if len(im.Data) > ImageDataMaxLen {
		return nil, &ImageError{Message: fmt.Sprintf("marshal: data length %d is greater than the maximum value %d", len(im.Data), ImageDataMaxLen)}
	}

	data := make([]byte, 0, 4+len(im.Format)+4+len(im.Data))
	buf := make([]byte, 4)

	imageByteOrder.PutUint32(buf, uint32(len(im.Format)))
	data = append(data, buf...)
	data = append(data, im.Format...)

	imageByteOrder.PutUint32(buf, uint32(len(im.Data)))
	data = append(data, buf...)
	data = append(data, im.Data...)

	return data, nil
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
	readData := func(length int) ([]byte, error) {
		if length > len(data) {
			return nil, &ImageError{Message: "unmarshal: unexpected end of data"}
		}
		res := data[:length]
		data = data[length:]
		return res, nil
	}

	var buf []byte
	var err error

	buf, err = readData(4)
	if err != nil {
		return err
	}
	formatLen := imageByteOrder.Uint32(buf)
	if formatLen > ImageFormatMaxLen {
		return &ImageError{Message: fmt.Sprintf("unmarshal: format length %d is greater than the maximum value %d", formatLen, ImageFormatMaxLen)}
	}

	buf, err = readData(int(formatLen))
	if err != nil {
		return err
	}
	im.Format = string(buf)

	buf, err = readData(4)
	if err != nil {
		return err
	}
	dataLen := imageByteOrder.Uint32(buf)
	if dataLen > ImageDataMaxLen {
		return &ImageError{Message: fmt.Sprintf("unmarshal: data length %d is greater than the maximum value %d", dataLen, ImageDataMaxLen)}
	}

	buf, err = readData(int(dataLen))
	if err != nil {
		return err
	}
	im.Data = buf

	return nil
}

// ImageError is an Image error.
type ImageError struct {
	Message string
}

func (err *ImageError) Error() string {
	return fmt.Sprintf("image error: %s", err.Message)
}
