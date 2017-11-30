package radixstore

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/armon/go-radix"
	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/testutil"
)

var sampleDataFile = "../../data/store/sample-radix.csv"

func TestNew(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	_, err := New(filein, dir, logger)

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
			r := csv.NewReader(strings.NewReader(c.input))
			err := buildTree(tree, r)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestLoad(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, dir, logger)

	type Case struct {
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		{dir, true, "loading a dir fails"},
		{sampleDataFile + ".hoge", true, "loading non-exising file fails"},
		{sampleDataFile, false, "loading existing file succeeds"},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
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
	s, _ := New(filein, dir, logger)
	_ = s.Load(sampleDataFile)

	type Data map[string]interface{}
	type Case struct {
		input         string
		expectError   bool
		expectedValue Data
		subtest       string
	}
	cases := []Case{
		{
			"aaa",
			true,
			nil,
			"non-existing key fails",
		},
		{
			"hoge",
			false,
			Data{
				"key":   "hoge",
				"value": "hoge!",
			},
			"exact match on key 'hoge' succeeds",
		},
		{
			"hogefuga",
			false,
			Data{
				"key":   "hoge",
				"value": "hoge!",
			},
			"prefix match on key on 'hogefuga' succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			v, err := s.Get([]byte(c.input))
			if c.expectError {
				assert.NotNil(t, err)
			} else {
				var d Data
				err = json.Unmarshal(v, &d)
				assert.Nil(t, err)
				assert.True(t, assert.ObjectsAreEqual(c.expectedValue, d))
			}
		})
	}
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, dir, logger)
	done := s.Start()

	file := filepath.Join(dir, "drop-in")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleDataFile)
		filein <- file
	}

	_, err := s.Get([]byte("hoge"))
	assert.Nil(t, err)

	close(filein)
	<-done
}
