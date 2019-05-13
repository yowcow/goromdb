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

// ErrorBucketNotFound bucket-not-found error type
type ErrorBucketNotFound struct {
	error
}

// ErrorKeyNotFound key-not-found error type
type ErrorKeyNotFound struct {
	error
}

// ErrorInternal internal error
type ErrorInternal struct {
	error
}

// BucketNotFoundError returns an error for bucket-not-found
func BucketNotFoundError(bucket []byte) error {
	return &ErrorBucketNotFound{
		fmt.Errorf("bucket not found error: %s", string(bucket)),
	}
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

// IsErrorBucketNotFound returns if it is an ErrorBucketNotFound
func IsErrorBucketNotFound(err error) bool {
	switch err.(type) {
	case ErrorBucketNotFound:
		return true
	case *ErrorBucketNotFound:
		return true
	default:
		return false
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
