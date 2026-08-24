package rule

import (
	"errors"

	"coldchain/internal/data"
	"coldchain/internal/store"
)

// Current returns the currently published rule, or ErrNotFound when none.
func (s *Service) Current() (*data.Rule, error) {
	var r data.Rule
	err := s.store.ReadJSON(s.currentPath(), &r)
	if errors.Is(err, store.ErrNotFound) {
		return nil, data.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}
