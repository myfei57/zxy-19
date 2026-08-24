package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ErrNotFound is returned when a requested file does not exist.
var ErrNotFound = errors.New("store: record not found")

// WriteJSON encodes value and durably writes it to path, overwriting any
// existing file. The file is synced before the call returns.
func (s *Store) WriteJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	_, writeErr := file.Write(data)
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

// ReadJSON decodes the record at path. A missing file maps to ErrNotFound.
func (s *Store) ReadJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("store: decode %s: %w", path, err)
	}
	return nil
}

// Exists reports whether a file exists at path.
func (s *Store) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// List returns file names in dir sorted lexically; a missing dir is empty.
func (s *Store) List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// Remove deletes path. Missing files are treated as already removed.
func (s *Store) Remove(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
