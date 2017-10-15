package protocol

import (
	"bufio"
	"fmt"
)

type Protocol interface {
	Parse([]byte) ([][]byte, error)
	Reply(*bufio.Writer, []byte, []byte)
	Finish(*bufio.Writer)
}

func InvalidCommandError(line []byte) error {
	return fmt.Errorf("invalid command: %s", string(line))
}
