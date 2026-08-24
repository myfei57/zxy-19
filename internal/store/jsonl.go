package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// AppendJSON appends one JSON line to a JSONL file, creating the file when
// missing. The append is synced before returning.
func (s *Store) AppendJSON(path string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(append(data, '\n'))
	syncErr := file.Sync()
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}

// ReadLines returns the raw lines of a JSONL file; a missing file is empty.
func (s *Store) ReadLines(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()
	var lines [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, append([]byte(nil), scanner.Bytes()...))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// DecodeLine unmarshals one JSON line into value.
func DecodeLine(line []byte, value any) error {
	if len(line) == 0 {
		return errors.New("store: empty jsonl line")
	}
	return json.Unmarshal(line, value)
}
