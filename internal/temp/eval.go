package temp

import "coldchain/internal/data"

// Evaluate compares one reading against a threshold and produces a verdict.
// Readings at or above the threshold are considered over temperature.
func (s *Service) Evaluate(reading data.TemperatureReading, threshold float64, version int) data.Verdict {
	overheat := reading.Celsius >= threshold
	reason := ""
	if overheat {
		reason = "temperature at or above the published threshold"
	}
	return data.Verdict{
		Overheat:  overheat,
		Reading:   reading.Celsius,
		Threshold: threshold,
		Version:   version,
		Reason:    reason,
	}
}
