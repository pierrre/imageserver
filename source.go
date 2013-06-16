package imageserver

type Source interface {
	Get(sourceId string) (image *Image, err error)
}
