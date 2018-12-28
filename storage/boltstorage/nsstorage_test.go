package boltstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleNSDBFile = "../../_data/store/sample-ns-boltdb.db"

func TestNewNS(t *testing.T) {
	NewNS()
}

func TestGetNS(t *testing.T) {
	s := NewNS()
	err := s.Load(sampleNSDBFile)

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
			[2][]byte{[]byte("ns1"), []byte("hoge")},
			[]byte("hoge1"),
			false,
		},
		{
			"namespace and key exists 2",
			[2][]byte{[]byte("ns2"), []byte("hoge")},
			[]byte("hoge2"),
			false,
		},
		{
			"namespace exists but key not exists",
			[2][]byte{[]byte("ns1"), []byte("fuga")},
			nil,
			true,
		},
		{
			"namespace not exists",
			[2][]byte{[]byte("ns3"), []byte("fuga")},
			nil,
			true,
		},
		{
			"namespace not specified",
			[2][]byte{nil, []byte("fuga")},
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
