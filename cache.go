package imageserver

// Image cache
type Cache interface {
	Get(key string, parameters Parameters) (image *Image, err error)
	Set(key string, image *Image, parameters Parameters) (err error)
}
