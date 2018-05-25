package storage

import (
	"fmt"
)

// Storage defines an interface to a storage
type Storage interface {
	Get(key []byte) ([]byte, error)
	Load(string) error
}

// NSStorage defines an interface to a namespaced storage
type NSStorage interface {
	Storage
	GetNS(namespace, key []byte) ([]byte, error)
}

// ErrorKeyNotFound key-not-found error type
type ErrorKeyNotFound struct {
	error
}

// ErrorInternal internal error
type ErrorInternal struct {
	error
}

// KeyNotFoundError returns an error for key-not-found
func KeyNotFoundError(key []byte) error {
	return &ErrorKeyNotFound{
		fmt.Errorf("key not found error: %s", string(key)),
	}
}

// InternalError something go wrong internal
func InternalError(s string) error {
	return &ErrorInternal{
		fmt.Errorf(s),
	}
}

// IsErrorKeyNotFound returns if it is an ErrorKeyNotFound
func IsErrorKeyNotFound(err error) bool {
	switch err.(type) {
	case ErrorKeyNotFound:
		return true
	case *ErrorKeyNotFound:
		return true
	default:
		return false
	}
}
