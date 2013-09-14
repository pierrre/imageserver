package imageserver

// Image cache
type Cache interface {
	Get(string, Parameters) (*Image, error)
	Set(string, *Image, Parameters) error
}
