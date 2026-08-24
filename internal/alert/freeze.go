package alert

import (
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Freeze handles an over-temperature reading for a batch and raises the
// frozen state together with the overheat evidence.
//
// The overheat record is persisted before the batch is frozen so that the
// freezing step always has durable evidence to lean on. Were the order
// reversed, a failure mid-way could leave a batch frozen with no overheat
// record on file, leaving the later thaw without grounds.
func (s *Service) Freeze(batchID string, reading data.TemperatureReading) (*data.Alert, error) {
	record := data.OverheatRecord{
		ID:        data.NewID(),
		BatchID:   batchID,
		Reading:   reading,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.temp.AppendOverheat(batchID, record); err != nil {
		return nil, err
	}
	if err := s.freezeState(batchID); err != nil {
		return nil, err
	}
	alertRecord := &data.Alert{
		ID:        data.NewID(),
		BatchID:   batchID,
		Severity:  "high",
		Message:   "batch frozen for over-temperature evidence",
		CreatedAt: time.Now().UTC(),
	}
	if err := store.Sanitize(alertRecord.ID); err != nil {
		return nil, err
	}
	if err := s.store.WriteJSON(s.path(alertRecord.ID), alertRecord); err != nil {
		return nil, err
	}
	return alertRecord, nil
}

// freezeState persists the frozen flag and moves the batch into the frozen
// status.
func (s *Service) freezeState(batchID string) error {
	flag := data.FrozenFlag{BatchID: batchID, Frozen: true, UpdatedAt: time.Now().UTC()}
	if err := s.store.WriteJSON(s.frozenFlagFile(batchID), flag); err != nil {
		return err
	}
	return s.batch.MarkFrozen(batchID)
}

// ClearFrozen durably removes the frozen flag of a batch.
func (s *Service) ClearFrozen(batchID string) error {
	if err := store.Sanitize(batchID); err != nil {
		return err
	}
	flag := data.FrozenFlag{BatchID: batchID, Frozen: false, UpdatedAt: time.Now().UTC()}
	return s.store.WriteJSON(s.frozenFlagFile(batchID), flag)
}
