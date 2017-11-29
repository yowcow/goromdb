package testutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTmpDir(t *testing.T) {
	dir := CreateTmpDir()
	defer os.RemoveAll(dir)

	_, err := os.Stat(dir)

	assert.Nil(t, err)
}

func TestCopyFile(t *testing.T) {
	dir := CreateTmpDir()
	defer os.RemoveAll(dir)

	dst := filepath.Join(dir, "hoge.txt")
	CopyFile(dst, "hoge.txt")

	fi, err := os.Open(dst)

	assert.Nil(t, err)

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fi)
	fi.Close()

	assert.Nil(t, err)
	assert.Equal(t, buf.String(), "hogehoge?\n")
}
