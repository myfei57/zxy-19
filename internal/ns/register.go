package ns

import (
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Register creates a new namespace record and persists it.
func (s *Service) Register(name, kind string, owner string) (*data.Namespace, error) {
	if name == "" || owner == "" {
		return nil, data.ErrInvalidInput
	}
	nsv := &data.Namespace{
		ID:        data.NewID(),
		Name:      name,
		Kind:      kind,
		Owner:     owner,
		CreatedAt: time.Now().UTC(),
	}
	if err := store.Sanitize(nsv.ID); err != nil {
		return nil, err
	}
	if err := s.store.WriteJSON(s.path(nsv.ID), nsv); err != nil {
		return nil, err
	}
	return nsv, nil
}
