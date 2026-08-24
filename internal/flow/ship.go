package flow

import (
	"time"

	"coldchain/internal/data"
)

// Ship moves a batch out of the warehouse. The dispatch notice is durably
// written first, then the batch state changes, and only a successful state
// change is reported in the audit trail.
func (s *Service) Ship(batchID string) (*data.Notice, error) {
	notice := &data.Notice{
		ID:        data.NewID(),
		BatchID:   batchID,
		Kind:      "dispatch",
		Message:   "batch dispatched from warehouse",
		CreatedAt: time.Now().UTC(),
	}
	if err := s.AppendNotice(batchID, *notice); err != nil {
		return nil, err
	}
	if err := s.batch.MarkShipped(batchID); err != nil {
		return nil, err
	}
	if err := s.audit.Success("ship", batchID, "batch shipped with durable notice"); err != nil {
		return nil, err
	}
	return notice, nil
}
