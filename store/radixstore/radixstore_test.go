package radixstore

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/armon/go-radix"
	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/reader"
	"github.com/yowcow/goromdb/testutil"
)

var sampleDataFile = "../../data/store/sample-data.csv"
var sampleDataFileGzipped = "../../data/store/sample-data.csv.gz"

func TestNew(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	_, err := New(filein, false, dir, nil, logger)

	assert.NotNil(t, err)

	_, err = New(filein, false, dir, reader.NewCSVReader, logger)

	assert.Nil(t, err)
}

func TestBuildTree(t *testing.T) {
	type Case struct {
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			"hoge,fuga,foo\nfuga,hoge,bar\n",
			true,
			"loading not 2-column CSV fails",
		},
		{
			"hoge,fuga\nfuga,hoge\n",
			false,
			"loading 2-column CSV succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			tree := radix.New()
			r := reader.NewCSVReader(strings.NewReader(c.input))
			err := buildTree(tree, r)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestLoad(t *testing.T) {
	type Case struct {
		input       string
		gzipped     bool
		expectError bool
		subtest     string
	}
	cases := []Case{
		{"/tmp", false, true, "loading a dir fails"},
		{sampleDataFile + ".hoge", false, true, "loading non-exising file fails"},
		{sampleDataFile, false, false, "loading existing file succeeds"},
		{sampleDataFileGzipped, true, false, "loading gzipped file succeeds"},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			dir := testutil.CreateTmpDir()
			defer os.RemoveAll(dir)

			filein := make(chan string)
			logbuf := new(bytes.Buffer)
			logger := log.New(logbuf, "", 0)
			s, _ := New(filein, c.gzipped, dir, reader.NewCSVReader, logger)

			err := s.Load(c.input)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, false, dir, reader.NewCSVReader, logger)
	_ = s.Load(sampleDataFile)

	type Case struct {
		input         string
		expectError   bool
		expectedKey   []byte
		expectedValue []byte
		subtest       string
	}
	cases := []Case{
		{
			"aaa",
			true,
			nil,
			nil,
			"non-existing key fails",
		},
		{
			"hoge",
			false,
			[]byte("hoge"),
			[]byte("hoge!"),
			"exact match on key 'hoge' succeeds",
		},
		{
			"hogefuga",
			false,
			[]byte("hoge"),
			[]byte("hoge!"),
			"prefix match on key on 'hogefuga' succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			k, v, err := s.Get([]byte(c.input))
			assert.Equal(t, c.expectError, err != nil)
			assert.Equal(t, c.expectedKey, k)
			assert.Equal(t, c.expectedValue, v)
		})
	}
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, false, dir, reader.NewCSVReader, logger)
	done := s.Start()

	file := filepath.Join(dir, "drop-in")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleDataFile)
		filein <- file
	}

	_, _, err := s.Get([]byte("hoge"))
	assert.Nil(t, err)

	close(filein)
	<-done
}
