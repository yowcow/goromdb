package boltstorage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleDBFile = "../../data/store/sample-boltdb.db"

func TestNew(t *testing.T) {
	New("mybucket")
}

func TestLoad(t *testing.T) {
	s := New("goromdb")

	type Case struct {
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		//{
		//	"./",
		//	true,
		//	"loading directory fails",
		//},
		//{
		//	sampleDBFile + ".hoge",
		//	true,
		//	"loading non-existing file fails",
		//},
		{
			sampleDBFile,
			false,
			"loading valid file succeeds",
		},
		{
			sampleDBFile,
			false,
			"loading valid file again succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			err := s.Load(sampleDBFile)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestLoadAndIterate(t *testing.T) {
	s := New("goromdb")

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

	err := s.LoadAndIterate(sampleDBFile, iterFunc)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(data))
}

func TestGet(t *testing.T) {
	s := New("goromdb")
	s.Load(sampleDBFile)

	type Case struct {
		input       []byte
		expectedVal []byte
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			[]byte("hogehoge"),
			nil,
			true,
			"non-existing key fails",
		},
		{
			[]byte("hoge"),
			[]byte("hoge!"),
			false,
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
