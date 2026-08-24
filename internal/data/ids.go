package data

import "github.com/google/uuid"

// NewID returns a fresh random identifier used for namespaces, probes,
// batches, readings, alerts, notices and audit events.
func NewID() string {
	return uuid.NewString()
}
