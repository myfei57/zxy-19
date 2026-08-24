package probe

import "coldchain/internal/data"

// Append stores one temperature reading for a probe under its reporting
// quota. Quota is checked before the reading is persisted so that a reading
// over the limit is rejected before it can occupy disk.
func (s *Service) Append(probeID string, reading data.TemperatureReading) error {
	if err := s.quota.Check(probeID); err != nil {
		return err
	}
	if err := s.temp.AppendReading(probeID, reading); err != nil {
		return err
	}
	return s.quota.Consume(probeID)
}
