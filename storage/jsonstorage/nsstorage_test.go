package jsonstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNS(t *testing.T) {
	s := NewNS(false)
	v, err := s.GetNS([]byte("ns"), []byte("key"))

	assert.Nil(t, v)
	assert.NotNil(t, err)
}

func TestGetNS(t *testing.T) {
	s := NewNS(false)
	err := s.Load("valid-ns.json")

	assert.Nil(t, err)

	type Case struct {
		subtest       string
		input         [2][]byte
		expectedVal   []byte
		errorExpected bool
	}
	cases := []Case{
		{
			"namespace and key exists",
			[2][]byte{[]byte("hoge"), []byte("fuga")},
			[]byte("hoge-fuga!"),
			false,
		},
		{
			"namespace exists but key not exists",
			[2][]byte{[]byte("hoge"), []byte("hoge")},
			nil,
			true,
		},
		{
			"namespace not exists",
			[2][]byte{[]byte("fuga"), []byte("fuga")},
			nil,
			true,
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			v, err := s.GetNS(c.input[0], c.input[1])

			assert.Equal(t, c.expectedVal, v)
			assert.True(t, c.errorExpected == (err != nil))
		})
	}
}
