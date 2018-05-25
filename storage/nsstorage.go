package storage

// NSStorage Storage that supports namespace
type NSStorage interface {
	Storage
	GetNS(namespace, key []byte) ([]byte, error)
}
