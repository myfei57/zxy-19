package probe

import (
	"path/filepath"

	"coldchain/internal/quota"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

const kind = "probes"

// Service manages probes, their durable cursor and their report path.
type Service struct {
	store *store.Store
	temp  *temp.Service
	quota *quota.Service
}

// New returns a probe service backed by the given collaborators.
func New(st *store.Store, temps *temp.Service, quotas *quota.Service) *Service {
	return &Service{store: st, temp: temps, quota: quotas}
}

// path returns the JSON file path of a probe record.
func (s *Service) path(id string) string {
	return s.store.Path(kind, id+".json")
}

// dir returns the directory holding probe records.
func (s *Service) dir() string {
	return filepath.Join(s.store.Root(), kind)
}
