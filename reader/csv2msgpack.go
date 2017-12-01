package reader

import (
	"encoding/csv"
	"io"

	"gopkg.in/vmihailenco/msgpack.v2"
)

// CSV2MsgpackReader represents a CSV reader that converts row into Msgpack data
type CSV2MsgpackReader struct {
	r    *csv.Reader
	cols []string
}

type msgpackRowData map[string]string

// NewCSV2MsgpackReader returns a CSV2MsgpackReader
func NewCSV2MsgpackReader(r io.Reader) Reader {
	csvr := csv.NewReader(r)
	cols, _ := csvr.Read()
	return &CSV2MsgpackReader{csvr, cols}
}

// Read reads a line of CSV, and returns key with Msgpack-formatted value
func (r CSV2MsgpackReader) Read() ([]byte, []byte, error) {
	rec, err := r.r.Read()
	if err != nil {
		return nil, nil, err
	}
	row := make(msgpackRowData)
	for i, col := range r.cols {
		row[col] = rec[i]
	}
	v, err := msgpack.Marshal(row)
	if err != nil {
		return nil, nil, err
	}
	return []byte(rec[0]), v, nil
}
