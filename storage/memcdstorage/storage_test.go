package memcdstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/storage/bdbstorage"
)

var sampleDBFile = "../../_data/store/sample-memcachedb-bdb.db"

func TestNew(t *testing.T) {
	p := bdbstorage.New()
	New(p)
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
			"loading existing file succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			p := bdbstorage.New()
			s := New(p)
			err := s.Load(c.input)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	p := bdbstorage.New()
	s := New(p)
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
			"existing key succeeds",
		},
		{
			[]byte("hoge"),
			[]byte("hoge!"),
			false,
			"existing key again succeeds",
		},
		{
			[]byte("hogefuga"),
			nil,
			true,
			"non-existing key fails",
		},
		{
			[]byte("hogefuga"),
			nil,
			true,
			"non-existing key again fails",
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
