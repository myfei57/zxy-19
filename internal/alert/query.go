package alert

import "coldchain/internal/data"

// List returns every alert sorted by file name.
func (s *Service) List() ([]data.Alert, error) {
	names, err := s.store.List(s.dir())
	if err != nil {
		return nil, err
	}
	out := make([]data.Alert, 0, len(names))
	for _, name := range names {
		id := name[:len(name)-len(".json")]
		var a data.Alert
		if err := s.store.ReadJSON(s.path(id), &a); err != nil {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

// Count returns the total number of alerts.
func (s *Service) Count() (int, error) {
	all, err := s.List()
	if err != nil {
		return 0, err
	}
	return len(all), nil
}
