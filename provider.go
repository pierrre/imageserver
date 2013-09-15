package imageserver

// Image provider interface
type Provider interface {
	Get(source interface{}, parameters Parameters) (*Image, error)
}
