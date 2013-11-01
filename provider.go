package imageserver

// Provider represents an Image provider
type Provider interface {
	Get(source interface{}, parameters Parameters) (*Image, error)
}
