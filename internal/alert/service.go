package alert

import (
	"coldchain/internal/batch"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// Service creates overheat alerts and freezes batches with durable evidence.
type Service struct {
	store *store.Store
	temp  *temp.Service
	batch *batch.Service
}

// New returns an alert service backed by the given collaborators.
func New(st *store.Store, temps *temp.Service, batches *batch.Service) *Service {
	return &Service{store: st, temp: temps, batch: batches}
}
