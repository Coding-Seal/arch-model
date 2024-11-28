package jsonl

import (
	"bufio"
	"encoding/json"
	"io"
)

type Scanner struct {
	sc *bufio.Scanner
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{sc: bufio.NewScanner(reader)}
}

func (sc *Scanner) Scan() bool {
	return sc.sc.Scan()
}

func (sc *Scanner) Err() error {
	return sc.sc.Err()
}

func (sc *Scanner) Json(v any) error {
	return json.Unmarshal(sc.sc.Bytes(), v)
}
