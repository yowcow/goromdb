package protocol

import (
	"fmt"
	"io"
)

type Protocol interface {
	Parse([]byte) ([][]byte, error)
	Reply(io.Writer, []byte, []byte)
	Finish(io.Writer)
}

func InvalidCommandError(line []byte) error {
	return fmt.Errorf("invalid command: %s", string(line))
}
