package bdbstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNS(t *testing.T) {
	NewNS()
}

func TestGetNS(t *testing.T) {
	s := NewNS()
	err := s.Load(sampleDBFile)

	assert.Nil(t, err)

	type Case struct {
		subtest       string
		input         [2][]byte
		expectedValue []byte
		errorExpected bool
	}
	cases := []Case{
		{
			"namespace and key exists",
			[2][]byte{[]byte("ho"), []byte("ge")},
			[]byte("hoge!"),
			false,
		},
		{
			"namespace exists but key not exists",
			[2][]byte{[]byte("ho"), []byte("ho")},
			nil,
			true,
		},
		{
			"namespace and key not exist",
			[2][]byte{[]byte("foo"), []byte("bar")},
			nil,
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			v, err := s.GetNS(c.input[0], c.input[1])

			assert.Equal(t, c.expectedValue, v)
			assert.True(t, c.errorExpected == (err != nil))
		})
	}
}
