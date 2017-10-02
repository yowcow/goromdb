package store

type Store interface {
	Get(string) (string, error)
}
