package fsitem

type FSItem interface {
	ToJson() (string, error)
	GetPath(childPath []byte) []byte
}
