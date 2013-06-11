package imageserver

type Source interface {
	Get(sourceId string) (*Image, error)
}
