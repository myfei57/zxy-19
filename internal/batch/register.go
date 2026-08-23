package batch

import (
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Register creates a batch in the received state and persists it.
func (s *Service) Register(code, namespaceID string) (*data.Batch, error) {
	if code == "" || namespaceID == "" {
		return nil, data.ErrInvalidInput
	}
	now := time.Now().UTC()
	b := &data.Batch{
		ID:          data.NewID(),
		Code:        code,
		NamespaceID: namespaceID,
		Status:      data.StatusReceived,
		Frozen:      false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := store.Sanitize(b.ID); err != nil {
		return nil, err
	}
	if err := s.store.WriteJSON(s.path(b.ID), b); err != nil {
		return nil, err
	}
	return b, nil
}
