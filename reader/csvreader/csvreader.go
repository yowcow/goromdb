package csvreader

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/yowcow/goromdb/reader"
)

type Reader struct {
	r *csv.Reader
}

func New(r io.Reader) reader.Reader {
	csvr := csv.NewReader(r)
	return &Reader{csvr}
}

func (r Reader) Read() ([]byte, []byte, error) {
	rec, err := r.r.Read()
	if err != nil {
		return nil, nil, err
	}
	if len(rec) != 2 {
		return nil, nil, fmt.Errorf("csvreader cannot read a row with a number of elements not exactly 2: got %d", len(rec))
	}
	return []byte(rec[0]), []byte(rec[1]), nil
}
