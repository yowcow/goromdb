package test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateStoreDir creates a temporary data store directory for testing
func CreateStoreDir() (string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "goromdb-test")
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// CopyDBFile copies content of dbfile into file data.db in given directory
func CopyDBFile(dir, dbfile string) (string, error) {
	file := filepath.Join(dir, "data.db")
	fw, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer fw.Close()

	fr, err := os.Open(dbfile)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	_, err = io.Copy(fw, fr)
	if err != nil {
		return "", err
	}

	md5file := filepath.Join(dir, "data.db.md5")
	fw, err = os.OpenFile(md5file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer fw.Close()

	fr, err = os.Open(dbfile + ".md5")
	if err != nil {
		return "", err
	}
	defer fr.Close()

	_, err = io.Copy(fw, fr)
	if err != nil {
		return "", err
	}

	return file, nil
}
