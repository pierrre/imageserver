package imageserver

type Processor interface {
	Process(inImage *Image, parameters Parameters) (image *Image, err error)
}
