package alert

import "path/filepath"

const kind = "alerts"

// path returns the JSON file path of an alert record.
func (s *Service) path(id string) string {
	return s.store.Path(kind, id+".json")
}

// dir returns the directory holding alert records.
func (s *Service) dir() string {
	return filepath.Join(s.store.Root(), kind)
}

// frozenFlagFile returns the durable frozen flag file of a batch.
func (s *Service) frozenFlagFile(batchID string) string {
	return s.store.Path("frozen", batchID+".json")
}
