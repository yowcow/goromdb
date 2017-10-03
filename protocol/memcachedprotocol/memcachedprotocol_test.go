package memcachedprotocol

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
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)

	p, _ := New()
	p.Reply(w, "hoge", "hogefuga")
	err := w.Flush()

	assert.Nil(t, err)
	assert.Equal(t, "VALUE hoge 0 8\r\nhogefuga\r\n", buf.String())
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
