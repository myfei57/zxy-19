package temp

import "coldchain/internal/store"

// Service owns temperature readings, aggregation windows and summaries.
type Service struct {
	store *store.Store
}

// New returns a temperature service backed by the given store.
func New(st *store.Store) *Service {
	return &Service{store: st}
}
