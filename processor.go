package imageserver

// Processor represents an Image processor
type Processor interface {
	Process(*Image, Parameters) (*Image, error)
}
