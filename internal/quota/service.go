package quota

import "coldchain/internal/store"

// Service tracks per-probe reporting quotas.
type Service struct {
	store *store.Store
}

// New returns a quota service backed by the given store.
func New(st *store.Store) *Service {
	return &Service{store: st}
}
