package probe

import (
	"time"

	"coldchain/internal/data"
)

// Send forwards the readings a probe has produced to the downstream sink.
// The sink acknowledgement is durable before the probe cursor advances, so a
// failed ack never skips readings the sink never received.
func (s *Service) Send(probeID string) error {
	readings, err := s.temp.ListReadings(probeID)
	if err != nil {
		return err
	}
	cursor, err := s.Cursor(probeID)
	if err != nil {
		return err
	}
	ack := data.SinkBatch{
		ID:        data.NewID(),
		ProbeID:   probeID,
		From:      cursor.Position + 1,
		To:        len(readings),
		Count:     len(readings) - cursor.Position,
		CreatedAt: time.Now().UTC(),
	}
	return s.AdvanceAfterAck(probeID, ack)
}
