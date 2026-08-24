package rule

import "coldchain/internal/data"

// Check evaluates one reading against the published rule version.
func (s *Service) Check(reading data.TemperatureReading) (data.Verdict, error) {
	if s.snapshot == nil {
		current, err := s.Current()
		if err != nil {
			return data.Verdict{}, err
		}
		s.snapshot = current
	}
	return s.temp.Evaluate(reading, s.snapshot.Threshold, s.snapshot.Version), nil
}
