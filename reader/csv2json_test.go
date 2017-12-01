package reader

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSV2JSONReader(t *testing.T) {
	r := strings.NewReader("")
	NewCSV2JSONReader(r)
}

func TestCSV2JSONReaderRead(t *testing.T) {
	r := strings.NewReader(`z,a
1,2
3,4
5,6
7
`)
	cjr := NewCSV2JSONReader(r)

	type Expected struct {
		key []byte
		val []byte
	}
	expected := []Expected{
		{
			[]byte("1"),
			[]byte(`{"a":"2","z":"1"}`),
		},
		{
			[]byte("3"),
			[]byte(`{"a":"4","z":"3"}`),
		},
		{
			[]byte("5"),
			[]byte(`{"a":"6","z":"5"}`),
		},
	}

	for _, e := range expected {
		k, v, err := cjr.Read()

		assert.Nil(t, err)
		assert.Equal(t, e.key, k)
		assert.Equal(t, e.val, v)
	}

	_, _, err := cjr.Read()

	assert.NotNil(t, err)
	assert.NotEqual(t, io.EOF, err)

	_, _, err = cjr.Read()

	assert.Equal(t, io.EOF, err)
}
