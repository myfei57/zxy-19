package trace

import "coldchain/internal/data"

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
	return s.temp.LoadWindows(batchID, generation)
}
