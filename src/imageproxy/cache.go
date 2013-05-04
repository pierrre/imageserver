package imageproxy

type Cache interface {
	Get(key string) *Image
	Set(key string, image *Image)
}
