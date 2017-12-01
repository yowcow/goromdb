package reader

import (
	"encoding/csv"
	"fmt"
	"io"
)

// CSVReader represents a simple CSV reader
type CSVReader struct {
	r *csv.Reader
}

// NewCSVReader creates a CSVReader
func NewCSVReader(r io.Reader) Reader {
	csvr := csv.NewReader(r)
	return &CSVReader{csvr}
}

// Read reads a line of CSV, and returns key-value pair
func (r CSVReader) Read() ([]byte, []byte, error) {
	rec, err := r.r.Read()
	if err != nil {
		return nil, nil, err
	}
	if len(rec) != 2 {
		return nil, nil, fmt.Errorf("csvreader cannot read a row with a number of elements not exactly 2: got %d", len(rec))
	}
	return []byte(rec[0]), []byte(rec[1]), nil
}
