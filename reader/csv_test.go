package reader

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVReader(t *testing.T) {
	rdr := strings.NewReader("")
	NewCSVReader(rdr)
}

func TestCSVReaderRead(t *testing.T) {
	type Case struct {
		input         string
		shouldSucceed bool
		subtest       string
	}
	cases := []Case{
		{
			"hoge\n",
			false,
			"1-column row fails",
		},
		{
			"hoge,fuga,foo\n",
			false,
			"3-column row fails",
		},
		{
			"hoge,fuga\n1,2,3\n",
			false,
			"inconsistend column count row fails",
		},
		{
			"hoge,fuga\n1,2\n",
			true,
			"2-column rows succeeds",
		},
	}

	returnReaderError := func(r Reader) error {
		for {
			_, _, err := r.Read()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
		}
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			rdr := strings.NewReader(c.input)
			r := NewCSVReader(rdr)
			err := returnReaderError(r)

			assert.Equal(t, c.shouldSucceed, err == nil)
		})
	}
}

func TestCSVReaderReadReturnsExpectedKeyValue(t *testing.T) {
	rdr := strings.NewReader(`hoge,hoge!
fuga,fuga!!
foo,foo!!!
bar,bar!!!!
`)
	r := NewCSVReader(rdr)

	type Expected struct {
		key []byte
		val []byte
	}
	expected := []Expected{
		{[]byte("hoge"), []byte("hoge!")},
		{[]byte("fuga"), []byte("fuga!!")},
		{[]byte("foo"), []byte("foo!!!")},
		{[]byte("bar"), []byte("bar!!!!")},
	}

	for _, exp := range expected {
		k, v, _ := r.Read()
		assert.Equal(t, exp.key, k)
		assert.Equal(t, exp.val, v)
	}

	_, _, err := r.Read()
	assert.Equal(t, io.EOF, err)
}
