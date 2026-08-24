package trace

import (
	"sort"

	"coldchain/internal/data"
)

// Timeline returns the latest generation's windows ordered by start time.
func (s *Service) Timeline(batchID string) ([]data.Window, error) {
	latest, err := s.Latest(batchID)
	if err != nil {
		return nil, err
	}
	windows, err := s.temp.LoadWindows(batchID, latest)
	if err != nil {
		return nil, err
	}
	sort.Slice(windows, func(i, j int) bool {
		return windows[i].Start.Before(windows[j].Start)
	})
	return windows, nil
}
