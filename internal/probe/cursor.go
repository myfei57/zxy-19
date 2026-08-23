package probe

import (
	"errors"
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// cursorFile returns the durable cursor file of a probe.
func (s *Service) cursorFile(probeID string) string {
	return s.store.Path("probes", "cursor", probeID+".json")
}

// Cursor returns the last durably acknowledged position of a probe.
func (s *Service) Cursor(probeID string) (data.CursorState, error) {
	var state data.CursorState
	err := s.store.ReadJSON(s.cursorFile(probeID), &state)
	if errors.Is(err, store.ErrNotFound) {
		return data.CursorState{ID: probeID, Position: 0}, nil
	}
	if err != nil {
		return data.CursorState{}, err
	}
	return state, nil
}

// AdvanceCursor durably persists a new probe position.
func (s *Service) AdvanceCursor(probeID string, position int) error {
	if err := store.Sanitize(probeID); err != nil {
		return err
	}
	state := data.CursorState{ID: probeID, Position: position, UpdatedAt: time.Now().UTC()}
	return s.store.WriteJSON(s.cursorFile(probeID), state)
}

// AdvanceAfterAck durably persists the sink acknowledgement before the probe
// cursor advances, so a failed ack never skips readings the sink never
// received.
func (s *Service) AdvanceAfterAck(probeID string, ack data.SinkBatch) error {
	if err := s.temp.AppendAck(probeID, ack); err != nil {
		return err
	}
	return s.AdvanceCursor(probeID, ack.To)
}
