package flow

// Release lets a frozen batch back into circulation and clears its frozen
// state.
func (s *Service) Release(batchID string) error {
	if err := s.batch.MarkReleased(batchID); err != nil {
		return err
	}
	if err := s.alert.ClearFrozen(batchID); err != nil {
		return err
	}
	return s.audit.Success("release", batchID, "batch released after frozen flag cleared")
}
