package temp

import (
	"coldchain/internal/data"
)

// LoadWindows returns the windows of one batch filtered to a single
// generation. The generation parameter is part of the trace replay contract:
// replay must never mix windows across generations.
func (s *Service) LoadWindows(batchID string, generation int) ([]data.Window, error) {
	all, err := s.LoadAllWindows(batchID)
	if err != nil {
		return nil, err
	}
	out := make([]data.Window, 0, len(all))
	for _, win := range all {
		if win.Generation == generation {
			out = append(out, win)
		}
	}
	return out, nil
}

// LoadAllWindows returns every window of a batch across all generations.
func (s *Service) LoadAllWindows(batchID string) ([]data.Window, error) {
	dir := s.store.Path("temp", "windows", batchID)
	names, err := s.store.List(dir)
	if err != nil {
		return nil, err
	}
	out := make([]data.Window, 0, len(names))
	for _, name := range names {
		windowID := name[:len(name)-len(".json")]
		win, err := s.LoadWindow(batchID, windowID)
		if err != nil {
			return nil, err
		}
		out = append(out, *win)
	}
	return out, nil
}

// GenerationCount returns how many generations produced windows for a batch.
func (s *Service) GenerationCount(batchID string) (int, error) {
	all, err := s.LoadAllWindows(batchID)
	if err != nil {
		return 0, err
	}
	seen := map[int]struct{}{}
	for _, win := range all {
		seen[win.Generation] = struct{}{}
	}
	return len(seen), nil
}
