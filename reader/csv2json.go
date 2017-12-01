package reader

import (
	"encoding/csv"
	"encoding/json"
	"io"
)

// CSV2JSONReader represents a CSV reader that converts a value into JSON
type CSV2JSONReader struct {
	r    *csv.Reader
	cols []string
}

type jsonRowData map[string]string

// NewCSV2JSONReader creates a CSV2JSONReader
func NewCSV2JSONReader(r io.Reader) Reader {
	csvr := csv.NewReader(r)
	cols, _ := csvr.Read()
	return &CSV2JSONReader{csvr, cols}
}

// Read reads a line of CSV, and returns key with JSON-formatted value
func (r CSV2JSONReader) Read() ([]byte, []byte, error) {
	rec, err := r.r.Read()
	if err != nil {
		return nil, nil, err
	}
	row := make(jsonRowData)
	for i, col := range r.cols {
		row[col] = rec[i]
	}
	v, err := json.Marshal(row)
	if err != nil {
		return nil, nil, err
	}
	return []byte(rec[0]), v, nil
}
