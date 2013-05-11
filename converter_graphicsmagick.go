package imageproxy

type GraphicsMagickConverter struct {
	Executable string
	TempDir    string
}

func (converter *GraphicsMagickConverter) Convert(sourceImage *Image, parameters *Parameters) (image *Image, err error) {
	return sourceImage, nil
}
