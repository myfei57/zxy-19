package batch

import (
	"path/filepath"

	"coldchain/internal/store"
)

const kind = "batches"

// Service manages batch registration and the batch state machine.
type Service struct {
	store *store.Store
}

// New returns a batch service backed by the given store.
func New(st *store.Store) *Service {
	return &Service{store: st}
}

// path returns the JSON file path of a batch record.
func (s *Service) path(id string) string {
	return s.store.Path(kind, id+".json")
}

// dir returns the directory holding batch records.
func (s *Service) dir() string {
	return filepath.Join(s.store.Root(), kind)
}
