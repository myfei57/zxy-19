package rule

import (
	"coldchain/internal/data"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// Service evaluates readings against published overheat thresholds.
type Service struct {
	store    *store.Store
	temp     *temp.Service
	snapshot *data.Rule
}

// New returns a rule service backed by the given store and temperature eval.
func New(st *store.Store, temps *temp.Service) *Service {
	return &Service{store: st, temp: temps}
}
