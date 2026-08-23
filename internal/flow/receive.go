package flow

// Receive completes a batch at the receiver and records the sign-off audit.
func (s *Service) Receive(batchID string) error {
	if err := s.batch.MarkSigned(batchID); err != nil {
		return err
	}
	return s.audit.Success("receive", batchID, "batch signed at receiver")
}
