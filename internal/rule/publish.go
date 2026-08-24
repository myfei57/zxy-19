package rule

import (
	"time"

	"coldchain/internal/data"
)

// Publish persists a new rule version and makes it the current threshold.
func (s *Service) Publish(name, metric string, threshold float64) (*data.Rule, error) {
	version, err := s.CurrentVersion()
	if err != nil {
		return nil, err
	}
	r := &data.Rule{
		ID:          data.NewID(),
		Version:     version + 1,
		Name:        name,
		Metric:      metric,
		Threshold:   threshold,
		PublishedAt: time.Now().UTC(),
	}
	if err := s.store.WriteJSON(s.path(r.Version), r); err != nil {
		return nil, err
	}
	if err := s.store.WriteJSON(s.currentPath(), r); err != nil {
		return nil, err
	}
	// Drop any cached snapshot so subsequent checks pick up the new threshold.
	s.snapshot = nil
	return r, nil
}
