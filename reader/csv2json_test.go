package reader

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSV2JSONReader(t *testing.T) {
	type Case struct {
		input, subtest string
	}
	cases := []Case{
		{
			"",
			"empty string succeeds",
		},
		{
			"key",
			"1-column string succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			r := strings.NewReader(c.input)
			NewCSV2JSONReader(r)
		})
	}
}

func TestCSV2JSONReaderRead(t *testing.T) {
	r := strings.NewReader(`key,z,a
item1,1,2
item2,3,4
item3,5,6
item4,7
`)
	cjr := NewCSV2JSONReader(r)

	type Expected struct {
		key []byte
		val []byte
	}
	expected := []Expected{
		{
			[]byte("item1"),
			[]byte(`{"a":"2","z":"1"}`),
		},
		{
			[]byte("item2"),
			[]byte(`{"a":"4","z":"3"}`),
		},
		{
			[]byte("item3"),
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
