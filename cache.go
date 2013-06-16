package imageserver

type Cache interface {
	Get(key string) (image *Image, err error)
	Set(key string, image *Image) error
}
