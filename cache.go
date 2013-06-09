package imageserver

type Cache interface {
	Get(key string) (*Image, error)
	Set(key string, image *Image) error
}
