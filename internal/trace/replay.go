package trace

import (
	"fmt"

	"coldchain/internal/data"
)

// Replay returns the temperature windows of one batch for a single
// generation. Generation 0 resolves to the latest generation.
func (s *Service) Replay(batchID string, generation int) ([]data.Window, error) {
	if generation == 0 {
		latest, err := s.Latest(batchID)
		if err != nil {
			return nil, err
		}
		generation = latest
	}
	all, err := s.temp.LoadAllWindows(batchID)
	if err != nil {
		return nil, err
	}
	bySlot := map[string]data.Window{}
	for _, win := range all {
		key := fmt.Sprintf("%d|%d", win.Start.Unix(), win.End.Unix())
		current, exists := bySlot[key]
		if !exists || win.Generation < current.Generation {
			bySlot[key] = win
		}
	}
	out := make([]data.Window, 0, len(bySlot))
	for _, win := range bySlot {
		out = append(out, win)
	}
	return out, nil
}
