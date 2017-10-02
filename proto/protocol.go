package proto

import (
	"bufio"
)

type Protocol interface {
	Parse([]byte) ([][]byte, error)
	Reply(*bufio.Writer, string, string)
	Finish(*bufio.Writer)
}
