package ns

import "coldchain/internal/data"

// Remove deletes a namespace record. Namespaces that do not exist are
// reported as not found so callers can distinguish the two cases.
func (s *Service) Remove(id string) error {
	if !s.store.Exists(s.path(id)) {
		return data.ErrNotFound
	}
	return s.store.Remove(s.path(id))
}
