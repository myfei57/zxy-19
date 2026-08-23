package batch

import (
	"errors"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Get loads one batch by id.
func (s *Service) Get(id string) (*data.Batch, error) {
	var b data.Batch
	err := s.store.ReadJSON(s.path(id), &b)
	if errors.Is(err, store.ErrNotFound) {
		return nil, data.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// List returns every batch sorted by file name.
func (s *Service) List() ([]data.Batch, error) {
	names, err := s.store.List(s.dir())
	if err != nil {
		return nil, err
	}
	out := make([]data.Batch, 0, len(names))
	for _, name := range names {
		id := name[:len(name)-len(".json")]
		b, err := s.Get(id)
		if err != nil {
			continue
		}
		out = append(out, *b)
	}
	return out, nil
}

// CountByStatus returns how many batches are currently in a status.
func (s *Service) CountByStatus(status string) (int, error) {
	all, err := s.List()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, b := range all {
		if b.Status == status {
			count++
		}
	}
	return count, nil
}
