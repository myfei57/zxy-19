package flow

import (
	"coldchain/internal/alert"
	"coldchain/internal/audit"
	"coldchain/internal/batch"
	"coldchain/internal/store"
)

// Service orchestrates shipment, release and receiver sign-off flows.
type Service struct {
	store *store.Store
	batch *batch.Service
	alert *alert.Service
	audit *audit.Service
}

// New returns a flow service backed by the given collaborators.
func New(st *store.Store, batches *batch.Service, alerts *alert.Service, audits *audit.Service) *Service {
	return &Service{store: st, batch: batches, alert: alerts, audit: audits}
}
