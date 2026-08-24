package batch

import (
	"time"

	"coldchain/internal/data"
)

// transition validates and persists one state machine step.
func (s *Service) transition(id, to string) error {
	b, err := s.Get(id)
	if err != nil {
		return err
	}
	if b.Frozen && to != data.StatusReleased {
		return data.ErrFrozen
	}
	if !data.ValidTransition(b.Status, to) {
		return data.ErrInvalidState
	}
	next := *b
	next.Status = to
	if to == data.StatusReleased {
		next.Frozen = false
	}
	next.UpdatedAt = time.Now().UTC()
	return s.store.WriteJSON(s.path(id), &next)
}
