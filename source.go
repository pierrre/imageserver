package imageserver

type Source interface {
	Get(sourceId string, parameters Parameters) (image *Image, err error)
}
