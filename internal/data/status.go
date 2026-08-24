package data

// Batch lifecycle statuses of the cold-chain state machine.
const (
	StatusReceived  = "received"
	StatusInStorage = "in_storage"
	StatusFrozen    = "frozen"
	StatusReleased  = "released"
	StatusShipped   = "shipped"
	StatusSigned    = "signed"
)

// ValidTransition reports whether the batch state machine allows moving from
// one status to another.
func ValidTransition(from, to string) bool {
	switch from {
	case StatusReceived:
		return to == StatusInStorage
	case StatusInStorage:
		return to == StatusFrozen || to == StatusShipped
	case StatusFrozen:
		return to == StatusReleased
	case StatusReleased:
		return to == StatusShipped
	case StatusShipped:
		return to == StatusSigned
	}
	return false
}

// Terminal reports whether a status ends the batch lifecycle.
func Terminal(status string) bool {
	return status == StatusSigned
}
