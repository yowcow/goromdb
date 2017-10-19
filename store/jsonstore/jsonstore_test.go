package jsonstore

import (
	"bytes"
	"log"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleDBFile = "../../data/store/sample-data.json"

func TestLoadJSON_returns_error_on_non_existing_file(t *testing.T) {
	_, err := LoadJSON("./hoge/fuga")

	assert.NotNil(t, err)
}

func TestLoadJSON_returns_error_on_invalid_JSON(t *testing.T) {
	_, err := LoadJSON("./jsonstore-invalid.json")

	assert.NotNil(t, err)
}

func TestNew(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New(sampleDBFile, logger)

	assert.NotNil(t, store)

	assert.Nil(t, store.Shutdown())
}

func TestNew_with_non_existing_file(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New("./jsonstore-hogefuga.json", logger)

	assert.NotNil(t, store)

	assert.Nil(t, store.Shutdown())

	re := regexp.MustCompile("no such file or directory")
	logline, err := buf.ReadString('\n')

	assert.Nil(t, err)
	assert.True(t, re.MatchString(logline))
}

func TestGet_on_existing_key(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New(sampleDBFile, logger)
	value, err := store.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, "hoge!", string(value))
	assert.Nil(t, store.Shutdown())
}

func TestGet_on_non_existing_key(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New(sampleDBFile, logger)
	value, err := store.Get([]byte("foobar"))

	assert.Nil(t, value)
	assert.NotNil(t, err)
	assert.Nil(t, store.Shutdown())
}
