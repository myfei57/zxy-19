package probe

import "coldchain/internal/data"

// Append stores one temperature reading for a probe only when the probe still
// has reporting quota. The quota gate runs before the durable write so an
// over-quota reading never consumes storage.
func (s *Service) Append(probeID string, reading data.TemperatureReading) error {
	if err := s.quota.Check(probeID); err != nil {
		return err
	}
	if err := s.temp.AppendReading(probeID, reading); err != nil {
		return err
	}
	return s.quota.Consume(probeID)
}
