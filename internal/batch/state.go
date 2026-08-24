package batch

import (
	"time"

	"coldchain/internal/data"
)

// MarkInStorage moves a received batch into storage.
func (s *Service) MarkInStorage(id string) error {
	return s.transition(id, data.StatusInStorage)
}

// MarkShipped moves a batch out of the warehouse.
func (s *Service) MarkShipped(id string) error {
	return s.transition(id, data.StatusShipped)
}

// MarkReleased moves a frozen batch back into circulation.
func (s *Service) MarkReleased(id string) error {
	return s.transition(id, data.StatusReleased)
}

// MarkSigned completes the batch lifecycle at the receiver.
func (s *Service) MarkSigned(id string) error {
	return s.transition(id, data.StatusSigned)
}

// MarkFrozen freezes a batch and durably persists the frozen marker.
func (s *Service) MarkFrozen(id string) error {
	b, err := s.Get(id)
	if err != nil {
		return err
	}
	if b.Frozen {
		return data.ErrFrozen
	}
	if !data.ValidTransition(b.Status, data.StatusFrozen) {
		return data.ErrInvalidState
	}
	next := *b
	next.Status = data.StatusFrozen
	next.Frozen = true
	next.UpdatedAt = time.Now().UTC()
	return s.store.WriteJSON(s.path(id), &next)
}
