package protocol

import (
	"fmt"
	"io"
)

// Protocol represents an interface for a protocol
type Protocol interface {
	Parse([]byte) ([][]byte, error)
	Reply(io.Writer, []byte, []byte)
	Finish(io.Writer)
}

// InvalidCommandError returns an error for invalid command line
func InvalidCommandError(line []byte) error {
	return fmt.Errorf("invalid command: %s", string(line))
}
