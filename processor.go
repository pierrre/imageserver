package imageserver

// Image processor interface
type Processor interface {
	Process(*Image, Parameters) (*Image, error)
}
