package jsonstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			s := New(c.gzipped)
			err := s.Load(c.input)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	s := New(false)
	s.Load("valid.json")

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
			[]byte("hogehoge"),
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
