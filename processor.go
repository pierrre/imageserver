package imageserver

// Processes an image and returns a new (or the same) image
type Processor interface {
	Process(*Image, Parameters) (*Image, error)
}
