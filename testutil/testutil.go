package testutil

import (
	"io"
	"io/ioutil"
	"os"
)

// CreateTmpDir creates a temp dir, and returns its path
func CreateTmpDir() string {
	dir, err := ioutil.TempDir(os.TempDir(), "goromdb-test")
	if err != nil {
		panic(err)
	}
	return dir
}

// CopyFile copies content of a source file into destination file
func CopyFile(dst, src string) {
	fo, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	fi, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	if _, err = io.Copy(fo, fi); err != nil {
		panic(err)
	}
}
