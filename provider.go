package imageserver

type Provider interface {
	Get(source interface{}, parameters Parameters) (image *Image, err error)
}
