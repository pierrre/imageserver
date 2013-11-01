package imageserver

// Processor represents an image processor
type Processor interface {
	Process(*Image, Parameters) (*Image, error)
}
