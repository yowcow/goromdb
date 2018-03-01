package storage

import (
	"fmt"
)

// IterationFunc defines an interface to a callback function for LoadAndIterate()
type IterationFunc func([]byte, []byte) error

// ErrorKeyNotFound key-not-found error type
type ErrorKeyNotFound struct {
	error
}

// ErrorInternal internal error
type ErrorInternal struct {
	error
}

// Storage defines an interface to a storage
type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string) error
	LoadAndIterate(string, IterationFunc) error
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
