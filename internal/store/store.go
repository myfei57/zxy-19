package store

import (
	"fmt"
	"os"
	"path/filepath"
)

// Store persists cold-chain records as JSON files under one data root.
type Store struct {
	root string
}

// Open creates the data root if needed and returns a Store rooted there.
func Open(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("store: create root %s: %w", root, err)
	}
	return &Store{root: root}, nil
}

// Root returns the data root of the store.
func (s *Store) Root() string {
	return s.root
}

// Path joins relative record paths below the data root.
func (s *Store) Path(parts ...string) string {
	return filepath.Join(append([]string{s.root}, parts...)...)
}
