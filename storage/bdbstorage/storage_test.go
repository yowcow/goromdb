package bdbstorage

import (
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

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			s := New()
			err := s.Load(c.input)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	s := New()
	v, err := s.Get([]byte("hoge"))

	assert.Nil(t, v)
	assert.NotNil(t, err)

	s.Load(sampleDBFile)

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
