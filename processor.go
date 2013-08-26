package imageserver

// Processes an image and returns a new (or the same) image
type Processor interface {
	Process(inImage *Image, parameters Parameters) (image *Image, err error)
}
