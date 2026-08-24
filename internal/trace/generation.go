package trace

// Latest returns the highest generation that produced windows for a batch.
func (s *Service) Latest(batchID string) (int, error) {
	all, err := s.temp.LoadAllWindows(batchID)
	if err != nil {
		return 0, err
	}
	latest := 0
	for _, win := range all {
		if win.Generation > latest {
			latest = win.Generation
		}
	}
	return latest, nil
}
