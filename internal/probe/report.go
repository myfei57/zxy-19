package probe

import (
	"time"

	"coldchain/internal/data"
)

// Report records a probe reading through the quota-gated append path and
// binds the reading to the probe's current batch.
func (s *Service) Report(probeID string, reading data.TemperatureReading) error {
	p, err := s.Get(probeID)
	if err != nil {
		return err
	}
	if reading.ID == "" {
		reading.ID = data.NewID()
	}
	if reading.RecordedAt.IsZero() {
		reading.RecordedAt = time.Now().UTC()
	}
	reading.ProbeID = probeID
	reading.BatchID = p.BatchID
	return s.Append(probeID, reading)
}
