package temp

import (
	"coldchain/internal/data"
	"coldchain/internal/store"
)

// AppendReading durably appends a temperature reading for a probe.
func (s *Service) AppendReading(probeID string, r data.TemperatureReading) error {
	if err := store.Sanitize(probeID); err != nil {
		return err
	}
	if r.ID == "" {
		r.ID = data.NewID()
	}
	return s.store.AppendJSON(s.readingsFile(probeID), r)
}

// AppendOverheat durably appends an overheat record for a batch.
func (s *Service) AppendOverheat(batchID string, rec data.OverheatRecord) error {
	if err := store.Sanitize(batchID); err != nil {
		return err
	}
	return s.store.AppendJSON(s.overheatFile(batchID), rec)
}

// AppendAck durably appends the sink acknowledgement of one forwarded batch.
func (s *Service) AppendAck(probeID string, ack data.SinkBatch) error {
	if err := store.Sanitize(probeID); err != nil {
		return err
	}
	return s.store.AppendJSON(s.ackFile(probeID), ack)
}

// AppendResult durably appends an aggregated window result for a batch.
func (s *Service) AppendResult(batchID string, result data.WindowResult) error {
	if err := store.Sanitize(batchID); err != nil {
		return err
	}
	return s.store.AppendJSON(s.resultFile(batchID), result)
}

// overheatFile returns the JSONL file holding a batch's overheat evidence.
func (s *Service) overheatFile(batchID string) string {
	return s.store.Path("temp", "overheat", batchID+".jsonl")
}

// ackFile returns the JSONL file holding a probe's sink acknowledgements.
func (s *Service) ackFile(probeID string) string {
	return s.store.Path("temp", "acks", probeID+".jsonl")
}

// resultFile returns the JSONL file holding a batch's window results.
func (s *Service) resultFile(batchID string) string {
	return s.store.Path("temp", "results", batchID+".jsonl")
}
