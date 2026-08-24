package rule

import "strconv"

// CurrentVersion returns the highest persisted rule version, or zero when no
// rule has ever been published.
func (s *Service) CurrentVersion() (int, error) {
	names, err := s.store.List(s.dir())
	if err != nil {
		return 0, err
	}
	latest := 0
	for _, name := range names {
		if len(name) < 6 || name[0] != 'v' || name[len(name)-5:] != ".json" {
			continue
		}
		n, err := strconv.Atoi(name[1 : len(name)-5])
		if err != nil {
			continue
		}
		if n > latest {
			latest = n
		}
	}
	return latest, nil
}
