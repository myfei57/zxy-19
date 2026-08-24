package rule

import (
	"path/filepath"
	"strconv"
)

const kind = "rules"

// path returns the JSON file path of a rule version.
func (s *Service) path(version int) string {
	return s.store.Path(kind, "v"+strconv.Itoa(version)+".json")
}

// currentPath returns the JSON file holding the current rule pointer.
func (s *Service) currentPath() string {
	return s.store.Path(kind, "current.json")
}

// dir returns the directory holding rule version files.
func (s *Service) dir() string {
	return filepath.Join(s.store.Root(), kind)
}
