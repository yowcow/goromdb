package gateway

// Gateway defines an interface to a store
type Gateway interface {
	Start() <-chan bool
	Load(string) error
	Get([]byte) ([]byte, []byte, error)
}
