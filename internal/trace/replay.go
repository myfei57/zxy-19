package trace

import (
	"sort"

	"coldchain/internal/data"
)

// Replay returns the temperature windows of one batch for a single
// generation. Generation 0 resolves to the latest generation. Replay must
// never mix windows across generations: a later re-aggregation of the same
// time slot is a separate record, not a patch over the earlier one, so
// loading by time slot and picking one generation would silently let an older
// generation overwrite a newer one. Only the records of the resolved
// generation are loaded, ordered by start time for deterministic display.
func (s *Service) Replay(batchID string, generation int) ([]data.Window, error) {
	if generation == 0 {
		latest, err := s.Latest(batchID)
		if err != nil {
			return nil, err
		}
		generation = latest
	}
	windows, err := s.temp.LoadWindows(batchID, generation)
	if err != nil {
		return nil, err
	}
	sort.Slice(windows, func(i, j int) bool {
		return windows[i].Start.Before(windows[j].Start)
	})
	return windows, nil
}
