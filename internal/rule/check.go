package rule

import "coldchain/internal/data"

// Check evaluates one reading against the currently published rule version.
// The current rule is re-read on every check so a threshold published after
// the service started is used immediately.
func (s *Service) Check(reading data.TemperatureReading) (data.Verdict, error) {
	current, err := s.Current()
	if err != nil {
		return data.Verdict{}, err
	}
	return s.temp.Evaluate(reading, current.Threshold, current.Version), nil
}
