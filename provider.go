package imageserver

// Provider is the interface that is used by Server to get a source image
type Provider interface {
	Get(source interface{}, parameters Parameters) (*Image, error)
}
