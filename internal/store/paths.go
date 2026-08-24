package store

import (
	"errors"
	"regexp"
)

var safeName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// Sanitize validates that a record id only contains safe path characters.
func Sanitize(id string) error {
	if !safeName.MatchString(id) {
		return errors.New("store: unsafe record id")
	}
	return nil
}
