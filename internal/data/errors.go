package data

import "errors"

var (
	// ErrNotFound is returned when a requested record does not exist.
	ErrNotFound = errors.New("coldchain: record not found")
	// ErrQuotaExceeded is returned when a probe has no quota left.
	ErrQuotaExceeded = errors.New("coldchain: probe quota exceeded")
	// ErrFrozen is returned when an operation requires an unfrozen batch.
	ErrFrozen = errors.New("coldchain: batch is frozen")
	// ErrInvalidState is returned for an illegal state machine transition.
	ErrInvalidState = errors.New("coldchain: invalid state transition")
	// ErrInvalidInput is returned for a malformed request payload.
	ErrInvalidInput = errors.New("coldchain: invalid input")
)
