package temp

import (
	"coldchain/internal/data"
	"coldchain/internal/store"
)

// AppendSummary durably appends one window summary for a batch.
func (s *Service) AppendSummary(batchID string, summary data.WindowSummary) error {
	if err := store.Sanitize(batchID); err != nil {
		return err
	}
	return s.store.AppendJSON(s.summaryFile(batchID), summary)
}

// Summaries returns every summary persisted for a batch.
func (s *Service) Summaries(batchID string) ([]data.WindowSummary, error) {
	lines, err := s.store.ReadLines(s.summaryFile(batchID))
	if err != nil {
		return nil, err
	}
	out := make([]data.WindowSummary, 0, len(lines))
	for _, line := range lines {
		var summary data.WindowSummary
		if err := store.DecodeLine(line, &summary); err != nil {
			return nil, err
		}
		out = append(out, summary)
	}
	return out, nil
}

// summaryFile returns the JSONL file holding a batch's window summaries.
func (s *Service) summaryFile(batchID string) string {
	return s.store.Path("temp", "summaries", batchID+".jsonl")
}
