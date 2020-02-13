package fsitem

type FSItem interface {
	ToJson() ([]byte, error)
	GetPath() []byte
}
