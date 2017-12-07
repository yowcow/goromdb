package jsonstorage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleDataFile = "../../data/store/sample-data.json"

func TestNew(t *testing.T) {
	s := New(false)
	v, err := s.Get([]byte("hoge"))

	assert.Nil(t, v)
	assert.NotNil(t, err)
}

func TestLoad(t *testing.T) {
	type Case struct {
		gzipped     bool
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			false,
			"./",
			true,
			"loading directory fails",
		},
		{
			false,
			"invalid.json",
			true,
			"loading invalid json fails",
		},
		{
			false,
			"valid.json",
			false,
			"loading valid json succeeds",
		},
		{
			true,
			"valid.json.gz",
			false,
			"loading valid gzipped json succeeds",
		},
	}

	mux := new(sync.RWMutex)
	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			s := New(c.gzipped)
			err := s.Load(c.input, mux)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	s := New(false)
	mux := new(sync.RWMutex)
	s.Load(sampleDataFile, mux)

	type Case struct {
		input       []byte
		expectError bool
		expectedVal []byte
		subtest     string
	}
	cases := []Case{
		{
			[]byte("foobar"),
			true,
			nil,
			"non-existing key fails",
		},
		{
			[]byte("hoge"),
			false,
			[]byte("hoge!"),
			"existing key succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			v, err := s.Get(c.input)
			assert.Equal(t, c.expectError, err != nil)
			assert.Equal(t, c.expectedVal, v)
		})
	}
}

func TestIterate(t *testing.T) {
	s := New(false)

	data := make(map[string]string)
	expected := [][]byte{
		[]byte("hoge"),
		[]byte("fuga"),
		[]byte("foo"),
		[]byte("bar"),
		[]byte("buz"),
	}
	iterFunc := func(k, v []byte) error {
		assert.Contains(t, expected, k)
		data[string(k)] = string(v)
		return nil
	}

	err := s.Iterate(iterFunc)

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(data))

	mux := new(sync.RWMutex)
	s.Load(sampleDataFile, mux)

	err = s.Iterate(iterFunc)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(data))
}
