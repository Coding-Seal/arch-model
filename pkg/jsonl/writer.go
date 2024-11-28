package jsonl

import (
	"bufio"
	"encoding/json"
	"io"
)

type Writer struct {
	wr *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{wr: bufio.NewWriter(w)}
}

func (w *Writer) Flush() error {
	return w.wr.Flush()
}

func (w *Writer) WriteJson(v any) error {
	marshal, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err = w.wr.Write(marshal); err != nil {
		return err
	}

	return w.wr.WriteByte('\n')
}
