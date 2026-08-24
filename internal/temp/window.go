package temp

import (
	"errors"
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Window status values persisted with each aggregation window.
const (
	WindowOpen   = "open"
	WindowClosed = "closed"
)

// cursorFile returns the durable cursor file of a batch.
func (s *Service) cursorFile(batchID string) string {
	return s.store.Path("temp", "cursor", batchID+".json")
}

// Cursor returns the last durably advanced window position of a batch.
func (s *Service) Cursor(batchID string) (data.CursorState, error) {
	var state data.CursorState
	err := s.store.ReadJSON(s.cursorFile(batchID), &state)
	if errors.Is(err, store.ErrNotFound) {
		return data.CursorState{ID: batchID, Position: 0}, nil
	}
	if err != nil {
		return data.CursorState{}, err
	}
	return state, nil
}

// AdvanceCursor durably persists a new window cursor position.
func (s *Service) AdvanceCursor(batchID string, position int) error {
	if err := store.Sanitize(batchID); err != nil {
		return err
	}
	state := data.CursorState{ID: batchID, Position: position, UpdatedAt: time.Now().UTC()}
	return s.store.WriteJSON(s.cursorFile(batchID), state)
}

// FinalizeWindow makes an aggregated window result durable before the window
// cursor advances, so a failed result write never skips the window.
func (s *Service) FinalizeWindow(batchID string, result data.WindowResult) error {
	if err := s.AppendResult(batchID, result); err != nil {
		return err
	}
	cursor, err := s.Cursor(batchID)
	if err != nil {
		return err
	}
	return s.AdvanceCursor(batchID, cursor.Position+1)
}

// CloseWindow finalizes an open window: the summary is persisted first and
// only then is the window durably marked closed.
func (s *Service) CloseWindow(batchID string, win data.Window) error {
	summary := data.WindowSummary{
		ID:           data.NewID(),
		BatchID:      batchID,
		Start:        win.Start,
		End:          win.End,
		Generation:   win.Generation,
		Records:      win.ReadingCount,
		OK:           win.ReadingCount > 0,
		SummarizedAt: time.Now().UTC(),
	}
	if err := s.AppendSummary(batchID, summary); err != nil {
		return err
	}
	win.Status = WindowClosed
	return s.SaveWindow(win)
}

// windowFile returns the JSON file holding one window record.
func (s *Service) windowFile(batchID, windowID string) string {
	return s.store.Path("temp", "windows", batchID, windowID+".json")
}

// SaveWindow durably persists a window record.
func (s *Service) SaveWindow(win data.Window) error {
	if err := store.Sanitize(win.BatchID); err != nil {
		return err
	}
	return s.store.WriteJSON(s.windowFile(win.BatchID, win.ID), win)
}

// LoadWindow reads one window record.
func (s *Service) LoadWindow(batchID, windowID string) (*data.Window, error) {
	var win data.Window
	err := s.store.ReadJSON(s.windowFile(batchID, windowID), &win)
	if errors.Is(err, store.ErrNotFound) {
		return nil, data.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &win, nil
}
