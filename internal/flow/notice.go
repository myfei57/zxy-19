package flow

import (
	"coldchain/internal/data"
	"coldchain/internal/store"
)

// noticeFile returns the JSONL file holding a batch's dispatch notices.
func (s *Service) noticeFile(batchID string) string {
	return s.store.Path("flow", "notices", batchID+".jsonl")
}

// AppendNotice durably appends a dispatch notice for a batch.
func (s *Service) AppendNotice(batchID string, notice data.Notice) error {
	if err := store.Sanitize(batchID); err != nil {
		return err
	}
	return s.store.AppendJSON(s.noticeFile(batchID), notice)
}

// Notices returns every notice persisted for a batch.
func (s *Service) Notices(batchID string) ([]data.Notice, error) {
	lines, err := s.store.ReadLines(s.noticeFile(batchID))
	if err != nil {
		return nil, err
	}
	out := make([]data.Notice, 0, len(lines))
	for _, line := range lines {
		var n data.Notice
		if err := store.DecodeLine(line, &n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}
