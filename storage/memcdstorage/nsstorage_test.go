package memcdstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/storage/bdbstorage"
)

func TestNewNS(t *testing.T) {
	p := bdbstorage.NewNS()
	NewNS(p)
}

func TestGetNS(t *testing.T) {
	p := bdbstorage.NewNS()
	s := NewNS(p)
	err := s.Load(sampleDBFile)

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
			[2][]byte{
				[]byte("ho"),
				[]byte("ge"),
			},
			[]byte("hoge!"),
			false,
		},
		{
			"namespace exists but key not exists",
			[2][]byte{
				[]byte("ho"),
				[]byte("ho"),
			},
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
