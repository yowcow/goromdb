package handler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/loader"
)

var (
	_ Handler   = (*testHandler)(nil)
	_ NSHandler = (*testNSHandler)(nil)
)

type testHandler struct {
	name string
}

func (h *testHandler) Get(k []byte) ([]byte, error) {
	return []byte(fmt.Sprintf("get %s from %s", string(k), h.name)), nil
}

func (h *testHandler) Load(file string) error {
	return nil
}

func (h *testHandler) Start(filein <-chan string, ldr *loader.Loader) <-chan bool {
	return nil
}

type testNSHandler struct {
	testHandler
}

func (h *testNSHandler) GetNS(ns, k []byte) ([]byte, error) {
	return []byte(fmt.Sprintf("get %s in ns %s from %s", string(k), string(ns), h.name)), nil
}

func TestRegisterHandler(t *testing.T) {
	m := NewMultiplexer()

	err := m.RegisterHandler("h1", new(testHandler))

	assert.Nil(t, err)

	err = m.RegisterHandler("h1", new(testHandler))

	assert.NotNil(t, err)
}

func TestRegisterNSHandler(t *testing.T) {
	m := NewMultiplexer()

	err := m.RegisterNSHandler("h1", new(testNSHandler))

	assert.Nil(t, err)

	err = m.RegisterNSHandler("h1", new(testNSHandler))

	assert.NotNil(t, err)
}

func TestGetHandler(t *testing.T) {
	m := NewMultiplexer()
	_ = m.RegisterHandler("h1", &testHandler{"h1"})
	hdr, err := m.GetHandler("h1")

	assert.NotNil(t, hdr)
	assert.Nil(t, err)

	hdr, err = m.GetHandler("h2")

	assert.Nil(t, hdr)
	assert.NotNil(t, err)
}

func TestGetNSHandler(t *testing.T) {
	m := NewMultiplexer()
	_ = m.RegisterNSHandler("h1", &testNSHandler{testHandler{"h1"}})
	hdr, err := m.GetNSHandler("h1")

	assert.NotNil(t, hdr)
	assert.Nil(t, err)

	hdr, err = m.GetNSHandler("h2")

	assert.Nil(t, hdr)
	assert.NotNil(t, err)
}
