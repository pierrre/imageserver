package imageserver

// Get an image from a source
type Provider interface {
	Get(source interface{}, parameters Parameters) (image *Image, err error)
}
