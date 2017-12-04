package radixgateway

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage/jsonstorage"
	"github.com/yowcow/goromdb/testutil"
)

var sampleDataFile = "../../data/store/sample-data.json"

func TestNew(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "radix.data")

	stg := jsonstorage.New(false)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	New(filein, ldr, stg, logger)
}

func TestLoadAndGet(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "radix.data")

	stg := jsonstorage.New(false)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	str := New(filein, ldr, stg, logger)
	err := str.Load(sampleDataFile)

	assert.Nil(t, err)

	type Case struct {
		input       []byte
		expectError bool
		expectedKey []byte
		expectedVal []byte
		subtest     string
	}
	cases := []Case{
		{
			[]byte("ho"),
			true,
			nil,
			nil,
			"non-existing key fails",
		},
		{
			[]byte("hoge"),
			false,
			[]byte("hoge"),
			[]byte("hoge!"),
			"exact match on key succeeds",
		},
		{
			[]byte("hogefuga"),
			false,
			[]byte("hoge"),
			[]byte("hoge!"),
			"prefix match on key succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			k, v, err := str.Get(c.input)
			assert.Equal(t, c.expectError, err != nil)
			assert.Equal(t, c.expectedKey, k)
			assert.Equal(t, c.expectedVal, v)
		})
	}
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "radix.data")

	stg := jsonstorage.New(false)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	str := New(filein, ldr, stg, logger)
	done := str.Start()

	file := filepath.Join(dir, "dropin.db")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleDataFile)
		filein <- file
	}

	key, val, err := str.Get([]byte("hogefuga"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge"), key)
	assert.Equal(t, []byte("hoge!"), val)

	close(filein)
	<-done
}
