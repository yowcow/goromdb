package reader

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/vmihailenco/msgpack.v2"
)

func TestNewCSV2MsgpackReader(t *testing.T) {
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
			NewCSV2MsgpackReader(r)
		})
	}
}

func TestCSV2MsgpackReaderReads(t *testing.T) {
	r := strings.NewReader(`key,x,a
item1,1,2
item2,3,4
item3,5,6
item4,7
`)
	cmr := NewCSV2MsgpackReader(r)

	type Expected struct {
		key []byte
		val msgpackRowData
	}
	expected := []Expected{
		{
			[]byte("item1"),
			msgpackRowData{
				"a": "2",
				"x": "1",
			},
		},
		{
			[]byte("item2"),
			msgpackRowData{
				"a": "4",
				"x": "3",
			},
		},
		{
			[]byte("item3"),
			msgpackRowData{
				"a": "6",
				"x": "5",
			},
		},
	}

	for _, e := range expected {
		k, v, err := cmr.Read()

		assert.Nil(t, err)
		assert.Equal(t, e.key, k)

		data := make(msgpackRowData)
		err = msgpack.Unmarshal(v, &data)

		assert.Nil(t, err)
		assert.True(t, reflect.DeepEqual(e.val, data))
	}

	_, _, err := cmr.Read()

	assert.NotNil(t, err)
	assert.NotEqual(t, io.EOF, err)

	_, _, err = cmr.Read()

	assert.Equal(t, io.EOF, err)
}
