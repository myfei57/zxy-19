package flow

// Release lets a frozen batch back into circulation. The frozen flag is
// durably cleared before the batch is marked released.
func (s *Service) Release(batchID string) error {
	if err := s.alert.ClearFrozen(batchID); err != nil {
		return err
	}
	if err := s.batch.MarkReleased(batchID); err != nil {
		return err
	}
	return s.audit.Success("release", batchID, "batch released after frozen flag cleared")
}
