package quota

import (
	"time"
)

// Recover resets a probe's used quota when the window rolls over, so an
// exhausted quota never leaks into the next window.
func (s *Service) Recover(probeID string) error {
	state, err := s.load(probeID)
	if err != nil {
		return err
	}
	current := windowKey(time.Now())
	if state.Window != current {
		state.Window = current
		state.Used = 0
		state.UpdatedAt = time.Now().UTC()
		return s.store.WriteJSON(s.path(probeID), state)
	}
	return nil
}
