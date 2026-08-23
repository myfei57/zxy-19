package quota

import (
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Consume marks one more reading against a probe's current window quota.
func (s *Service) Consume(probeID string) error {
	state, err := s.load(probeID)
	if err != nil {
		return err
	}
	state.Used++
	state.UpdatedAt = time.Now().UTC()
	if err := store.Sanitize(probeID); err != nil {
		return err
	}
	return s.store.WriteJSON(s.path(probeID), state)
}

// SetLimit configures a probe's reporting quota for each window.
func (s *Service) SetLimit(probeID string, limit int) error {
	if limit < 0 {
		return data.ErrInvalidInput
	}
	state, err := s.load(probeID)
	if err != nil {
		return err
	}
	state.Limit = limit
	state.UpdatedAt = time.Now().UTC()
	return s.store.WriteJSON(s.path(probeID), state)
}
