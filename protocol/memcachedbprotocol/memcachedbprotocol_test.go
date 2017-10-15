package memcachedbprotocol

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New()

	assert.Nil(t, err)
}

func TestParse_on_get_command(t *testing.T) {
	p, _ := New()
	words, err := p.Parse([]byte("get hoge"))

	assert.Nil(t, err)
	assert.Equal(t, 1, len(words))
	assert.Equal(t, []byte("hoge"), words[0])
}

func TestParse_on_gets_command(t *testing.T) {
	p, _ := New()
	words, err := p.Parse([]byte("gets hoge fuga"))

	assert.Nil(t, err)
	assert.Equal(t, 2, len(words))
	assert.Equal(t, []byte("hoge"), words[0])
	assert.Equal(t, []byte("fuga"), words[1])
}

func TestParse_on_invalid_command(t *testing.T) {
	p, _ := New()
	words, err := p.Parse([]byte("set hoge fuga foo bar"))

	assert.Equal(t, "invalid command: set hoge fuga foo bar", err.Error())
	assert.Equal(t, 0, len(words))
}

func TestReply(t *testing.T) {
	memdb := new(bytes.Buffer)
	err := Serialize(memdb, []byte("hoge"), []byte("hogefuga!!!"))

	assert.Nil(t, err)

	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)

	p, _ := New()
	p.Reply(w, []byte("hoge"), memdb.Bytes())
	err = w.Flush()

	assert.Nil(t, err)
	assert.Equal(t, "VALUE hoge 0 11\r\nhogefuga!!!\r\n", buf.String())
}

func TestFinish(t *testing.T) {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)

	p, _ := New()
	p.Finish(w)
	err := w.Flush()

	assert.Nil(t, err)
	assert.Equal(t, "END\r\n", buf.String())
}

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
