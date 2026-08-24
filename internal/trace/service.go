package trace

import (
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// Service replays batch temperature history by generation.
type Service struct {
	store *store.Store
	temp  *temp.Service
}

// New returns a trace service backed by the given collaborators.
func New(st *store.Store, temps *temp.Service) *Service {
	return &Service{store: st, temp: temps}
}
