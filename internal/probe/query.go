package probe

import (
	"errors"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Get loads one probe by id.
func (s *Service) Get(id string) (*data.Probe, error) {
	var p data.Probe
	err := s.store.ReadJSON(s.path(id), &p)
	if errors.Is(err, store.ErrNotFound) {
		return nil, data.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List returns every probe sorted by file name.
func (s *Service) List() ([]data.Probe, error) {
	names, err := s.store.List(s.dir())
	if err != nil {
		return nil, err
	}
	out := make([]data.Probe, 0, len(names))
	for _, name := range names {
		id := name[:len(name)-len(".json")]
		p, err := s.Get(id)
		if err != nil {
			continue
		}
		out = append(out, *p)
	}
	return out, nil
}

// ListByNamespace returns the probes attached to one namespace.
func (s *Service) ListByNamespace(namespaceID string) ([]data.Probe, error) {
	all, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]data.Probe, 0, len(all))
	for _, p := range all {
		if p.NamespaceID == namespaceID {
			out = append(out, p)
		}
	}
	return out, nil
}
