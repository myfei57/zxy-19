package ns

import (
	"errors"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Get loads one namespace by id.
func (s *Service) Get(id string) (*data.Namespace, error) {
	var nsv data.Namespace
	err := s.store.ReadJSON(s.path(id), &nsv)
	if errors.Is(err, store.ErrNotFound) {
		return nil, data.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &nsv, nil
}

// List returns every namespace sorted by file name.
func (s *Service) List() ([]data.Namespace, error) {
	names, err := s.store.List(s.dir())
	if err != nil {
		return nil, err
	}
	out := make([]data.Namespace, 0, len(names))
	for _, name := range names {
		id := name[:len(name)-len(".json")]
		nsv, err := s.Get(id)
		if err != nil {
			continue
		}
		out = append(out, *nsv)
	}
	return out, nil
}
