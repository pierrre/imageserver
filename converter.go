package imageproxy

type Converter interface {
	Convert(image *Image, parameters *Parameters) (*Image, error)
}
