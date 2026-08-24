package quota

import (
	"errors"
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Check verifies that a probe still has quota left in the current window.
func (s *Service) Check(probeID string) error {
	if err := s.Recover(probeID); err != nil {
		return err
	}
	state, err := s.load(probeID)
	if err != nil {
		return err
	}
	if state.Used >= state.Limit {
		return data.ErrQuotaExceeded
	}
	return nil
}

// load reads the durable quota state of a probe, defaulting to an empty
// state when the probe has never been configured.
func (s *Service) load(probeID string) (data.QuotaState, error) {
	var state data.QuotaState
	err := s.store.ReadJSON(s.path(probeID), &state)
	if errors.Is(err, store.ErrNotFound) {
		return data.QuotaState{
			ProbeID: probeID,
			Limit:   0,
			Used:    0,
			Window:  windowKey(time.Now()),
		}, nil
	}
	if err != nil {
		return data.QuotaState{}, err
	}
	return state, nil
}
