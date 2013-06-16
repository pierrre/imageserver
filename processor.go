package imageserver

type Processor interface {
	Process(image *Image, parameters Parameters) (*Image, error)
}
