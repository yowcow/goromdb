package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidCommandError(t *testing.T) {
	result := InvalidCommandError([]byte("hoge fuga"))

	assert.Equal(t, "invalid command: hoge fuga", result.Error())
}
