package audit

import (
	"time"

	"coldchain/internal/data"
)

// Record durably appends one audit event.
func (s *Service) Record(event data.AuditEvent) error {
	if event.ID == "" {
		event.ID = data.NewID()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	return s.store.AppendJSON(s.auditFile(), event)
}

// Success records a successful operation event for an actor and subject.
func (s *Service) Success(actor, subject, detail string) error {
	return s.Record(data.AuditEvent{
		ID:        data.NewID(),
		Actor:     actor,
		Action:    "complete",
		Subject:   subject,
		OK:        true,
		Detail:    detail,
		CreatedAt: time.Now().UTC(),
	})
}
