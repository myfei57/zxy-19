package probe

import (
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Register creates a probe attached to a namespace and optional batch.
func (s *Service) Register(namespaceID, name, batchID string) (*data.Probe, error) {
	if namespaceID == "" || name == "" {
		return nil, data.ErrInvalidInput
	}
	existing, err := s.ListByNamespace(namespaceID)
	if err != nil {
		return nil, err
	}
	p := &data.Probe{
		ID:          data.NewID(),
		Name:        name,
		NamespaceID: namespaceID,
		BatchID:     batchID,
		Sequence:    len(existing) + 1,
		CreatedAt:   time.Now().UTC(),
	}
	if err := store.Sanitize(p.ID); err != nil {
		return nil, err
	}
	if err := s.store.WriteJSON(s.path(p.ID), p); err != nil {
		return nil, err
	}
	return p, nil
}
