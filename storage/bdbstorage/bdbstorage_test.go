package bdbstorage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleDBFile = "../../data/store/sample-bdb.db"

func TestNew(t *testing.T) {
	New()
}

func TestLoad(t *testing.T) {
	type Case struct {
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			"./",
			true,
			"loading directory fails",
		},
		{
			sampleDBFile + ".hoge",
			true,
			"loading non-existing file fails",
		},
		{
			sampleDBFile,
			false,
			"loading valid file succeeds",
		},
	}

	mux := new(sync.RWMutex)
	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			s := New()
			err := s.Load(c.input, mux)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	s := New()
	mux := new(sync.RWMutex)
	s.Load(sampleDBFile, mux)

	type Case struct {
		input       []byte
		expected    []byte
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			[]byte("hoge"),
			[]byte("hoge!"),
			false,
			"existing key returns expected val",
		},
		{
			[]byte("hogehoge"),
			nil,
			true,
			"non-existing key returns error",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			v, err := s.Get(c.input)
			assert.Equal(t, c.expectError, err != nil)
			assert.Equal(t, c.expected, v)
		})
	}
}

func TestCursor(t *testing.T) {
	s := New()

	_, err := s.Cursor()

	assert.NotNil(t, err)

	mux := new(sync.RWMutex)
	s.Load(sampleDBFile, mux)
	c, err := s.Cursor()

	assert.Nil(t, err)

	expected := [][]byte{
		[]byte("hoge"),
		[]byte("fuga"),
		[]byte("foo"),
		[]byte("bar"),
		[]byte("buz"),
	}
	count := 0
	for {
		k, _, err := c.Next()
		if err != nil {
			break
		}
		assert.Contains(t, expected, k)
		count++
	}

	assert.Nil(t, c.Close())
	assert.Equal(t, len(expected), count)
}
