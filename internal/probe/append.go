package probe

import "coldchain/internal/data"

// Append stores one temperature reading for a probe under its reporting
// quota.
func (s *Service) Append(probeID string, reading data.TemperatureReading) error {
	if err := s.temp.AppendReading(probeID, reading); err != nil {
		return err
	}
	if err := s.quota.Check(probeID); err != nil {
		return err
	}
	return s.quota.Consume(probeID)
}
