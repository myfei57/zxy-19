package audit

import "coldchain/internal/store"

// Service records immutable operation audit events.
type Service struct {
	store *store.Store
}

// New returns an audit service backed by the given store.
func New(st *store.Store) *Service {
	return &Service{store: st}
}
