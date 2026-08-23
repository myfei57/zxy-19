package quota

import (
	"time"

	"coldchain/internal/data"
)

const kind = "quotas"

// path returns the JSON file path of a probe quota state.
func (s *Service) path(probeID string) string {
	return s.store.Path(kind, probeID+".json")
}

// windowKey returns the current quota window label.
func windowKey(now time.Time) string {
	return now.UTC().Format("2006-01-02-15")
}

// State returns the current quota state of a probe.
func (s *Service) State(probeID string) (data.QuotaState, error) {
	return s.load(probeID)
}
