package store

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Store represents an interface for a store
type Store interface {
	Get([]byte) ([]byte, error)
	Shutdown() error
}

// ChecksumChecker represents a function for checksum checking
type ChecksumChecker func(string) error

// CheckMD5Sum checks if given file's checksum file contains a valid MD5 sum
func CheckMD5Sum(file string) error {
	md5fh, err := os.Open(file + ".md5")
	if err != nil {
		return err
	}

	fh, err := os.Open(file)
	if err != nil {
		return err
	}

	defer func() {
		fh.Close()
		md5fh.Close()
	}()

	expected := make([]byte, 32)
	_, err = md5fh.Read(expected)
	if err != nil {
		return err
	}

	h := md5.New()
	if _, err := io.Copy(h, fh); err != nil {
		return err
	}

	md5sum := hex.EncodeToString(h.Sum(nil))
	if md5sum != string(expected) {
		return fmt.Errorf("expecting MD5 sum '%s' but got '%s'", expected, md5sum)
	}

	return nil
}

// KeyNotFoundError returns an error for key is not found
func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key '%s' not found", string(key))
}
