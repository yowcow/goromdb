package memcachedb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize_and_Deserialize(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Serialize(buf, []byte("hoge"), []byte("ほげほげ!!"))

	assert.Nil(t, err)

	r := bytes.NewReader(buf.Bytes())
	key, val, len, err := Deserialize(r)

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge"), key)
	assert.Equal(t, []byte("ほげほげ!!"), val)
	assert.Equal(t, 14, len)
}

func TestDeserialize_returns_header_error(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	_, _, _, err := Deserialize(r)

	assert.Equal(t, "failed reading memcachedb binary headers: EOF", err.Error())
}

func TestDeserialize_returns_body_error(t *testing.T) {
	r := bytes.NewReader([]byte("hogefuga"))
	_, _, _, err := Deserialize(r)

	assert.Equal(t, "failed reading memcachedb binary body: EOF", err.Error())
}
