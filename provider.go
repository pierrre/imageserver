package imageserver

type Provider interface {
	Get(source string, parameters Parameters) (image *Image, err error)
}
