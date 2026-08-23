package audit

import (
	"coldchain/internal/data"
	"coldchain/internal/store"
)

// List returns every audit event in append order.
func (s *Service) List() ([]data.AuditEvent, error) {
	lines, err := s.store.ReadLines(s.auditFile())
	if err != nil {
		return nil, err
	}
	out := make([]data.AuditEvent, 0, len(lines))
	for _, line := range lines {
		var event data.AuditEvent
		if err := store.DecodeLine(line, &event); err != nil {
			return nil, err
		}
		out = append(out, event)
	}
	return out, nil
}

// ListBySubject returns audit events mentioning one batch or probe.
func (s *Service) ListBySubject(subject string) ([]data.AuditEvent, error) {
	all, err := s.List()
	if err != nil {
		return nil, err
	}
	out := make([]data.AuditEvent, 0, len(all))
	for _, event := range all {
		if event.Subject == subject {
			out = append(out, event)
		}
	}
	return out, nil
}

// Count returns the total number of audit events.
func (s *Service) Count() (int, error) {
	all, err := s.List()
	if err != nil {
		return 0, err
	}
	return len(all), nil
}
