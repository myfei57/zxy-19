package temp

import (
	"time"

	"coldchain/internal/data"
)

// Aggregate computes a window result from persisted readings. The window
// result is written durably before the window cursor advances, so a failed
// result write never skips the window.
func (s *Service) Aggregate(batchID string, win data.Window) (*data.WindowResult, error) {
	if win.Status == "" {
		win.Status = WindowOpen
	}
	readings, err := s.ListReadings(win.ProbeID)
	if err != nil {
		return nil, err
	}
	result := data.WindowResult{
		ID:         data.NewID(),
		BatchID:    batchID,
		ProbeID:    win.ProbeID,
		Start:      win.Start,
		End:        win.End,
		Generation: win.Generation,
		ComputedAt: time.Now().UTC(),
	}
	for _, r := range readings {
		if r.RecordedAt.Before(win.Start) || !r.RecordedAt.Before(win.End) {
			continue
		}
		if result.Count == 0 {
			result.Min = r.Celsius
			result.Max = r.Celsius
		}
		if r.Celsius < result.Min {
			result.Min = r.Celsius
		}
		if r.Celsius > result.Max {
			result.Max = r.Celsius
		}
		result.Avg += r.Celsius
		result.Count++
	}
	if result.Count > 0 {
		result.Avg /= float64(result.Count)
	}
	return &result, s.FinalizeWindow(batchID, result)
}
