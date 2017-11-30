package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyNotFoundError(t *testing.T) {
	err := KeyNotFoundError([]byte("hogefuga"))

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "key 'hogefuga' not found")
}

func TestNewReader(t *testing.T) {
	type Case struct {
		file     string
		gzipped  bool
		expected string
		subtest  string
	}
	cases := []Case{
		{
			"test.txt",
			false,
			"hogehoge",
			"reading plain text succeeds",
		},
		{
			"test.txt.gz",
			true,
			"hogehoge",
			"reading gzipped text succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			f, _ := os.Open(c.file)
			r, err := NewReader(f, c.gzipped)
			assert.Nil(t, err)

			buf := make([]byte, 8)
			len, err := r.Read(buf)
			assert.Equal(t, 8, len)
			assert.Nil(t, err)
			assert.Equal(t, c.expected, string(buf))
		})
	}
}
