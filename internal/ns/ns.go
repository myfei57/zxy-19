package ns

import (
	"path/filepath"

	"coldchain/internal/store"
)

const kind = "namespaces"

// Service manages warehouse and carrier namespaces.
type Service struct {
	store *store.Store
}

// New returns a namespace service backed by the given store.
func New(st *store.Store) *Service {
	return &Service{store: st}
}

// path returns the JSON file path of a namespace record.
func (s *Service) path(id string) string {
	return s.store.Path(kind, id+".json")
}

// dir returns the directory holding namespace records.
func (s *Service) dir() string {
	return filepath.Join(s.store.Root(), kind)
}
