package audit

const kind = "audit"

// auditFile returns the JSONL file holding all audit events.
func (s *Service) auditFile() string {
	return s.store.Path(kind, "events.jsonl")
}
