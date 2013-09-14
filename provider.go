package imageserver

// Get an image from a source
type Provider interface {
	Get(interface{}, Parameters) (*Image, error)
}
