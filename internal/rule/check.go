package rule

import "coldchain/internal/data"

// Check evaluates one reading against the currently effective rule version.
// It always reads the live current rule so a newly published threshold takes
// effect immediately instead of being masked by a stale snapshot taken before
// the rule was republished.
func (s *Service) Check(reading data.TemperatureReading) (data.Verdict, error) {
	current, err := s.Current()
	if err != nil {
		return data.Verdict{}, err
	}
	s.snapshot = current
	return s.temp.Evaluate(reading, current.Threshold, current.Version), nil
}
