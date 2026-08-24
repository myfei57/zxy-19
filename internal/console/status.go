package console

import (
	"net/http"

	"coldchain/internal/data"
)

// status reports the current counts across the platform.
func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	namespaces, err := s.ns.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	probes, err := s.probe.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	batches, err := s.batch.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	alertCount, err := s.alert.Count()
	if err != nil {
		s.writeError(w, err)
		return
	}
	frozenCount, err := s.batch.CountByStatus(data.StatusFrozen)
	if err != nil {
		s.writeError(w, err)
		return
	}
	auditCount, err := s.audit.Count()
	if err != nil {
		s.writeError(w, err)
		return
	}
	completed := 0
	for _, b := range batches {
		if data.Terminal(b.Status) {
			completed++
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"namespaces":        len(namespaces),
		"probes":            len(probes),
		"batches":           len(batches),
		"alerts":            alertCount,
		"frozen_batches":    frozenCount,
		"audit_events":      auditCount,
		"completed_batches": completed,
	})
}
