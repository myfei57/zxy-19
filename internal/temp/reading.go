package temp

import (
	"path/filepath"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// readingsFile returns the JSONL file holding one probe's readings.
func (s *Service) readingsFile(probeID string) string {
	return s.store.Path("temp", "readings", probeID+".jsonl")
}

// readingsDir returns the directory holding per-probe reading files.
func (s *Service) readingsDir() string {
	return filepath.Join(s.store.Root(), "temp", "readings")
}

// ReadingCount returns how many readings a probe has persisted.
func (s *Service) ReadingCount(probeID string) (int, error) {
	lines, err := s.store.ReadLines(s.readingsFile(probeID))
	if err != nil {
		return 0, err
	}
	return len(lines), nil
}

// ListReadings decodes every reading of a probe in file order.
func (s *Service) ListReadings(probeID string) ([]data.TemperatureReading, error) {
	lines, err := s.store.ReadLines(s.readingsFile(probeID))
	if err != nil {
		return nil, err
	}
	out := make([]data.TemperatureReading, 0, len(lines))
	for _, line := range lines {
		var r data.TemperatureReading
		if err := store.DecodeLine(line, &r); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

// ListAllReadings returns every reading of every probe, used by the console
// temperature page when no probe filter is supplied.
func (s *Service) ListAllReadings() ([]data.TemperatureReading, error) {
	names, err := s.store.List(s.readingsDir())
	if err != nil {
		return nil, err
	}
	var out []data.TemperatureReading
	for _, name := range names {
		probeID := name[:len(name)-len(".jsonl")]
		items, err := s.ListReadings(probeID)
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
	}
	return out, nil
}
