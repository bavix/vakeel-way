package entities

// Status represents the status of a service.
type Status uint8

// String returns the string representation of the status.
//
// It returns "up" if the status is Up, "down" if the status is Down,
// and "Undefined" for any other value.
//
// Parameters:
//   - s: The Status value to convert to a string.
//
// Returns:
//   - A string representation of the status.
func (s Status) String() string {
	// Check the status and return the corresponding string.
	switch s {
	case Up:
		// The status is Up, so return "up".
		return "up"
	case Down:
		// The status is Down, so return "down".
		return "down"
	default:
		// The status is undefined, so return "Undefined".
		return "Undefined"
	}
}

// Status constants represent different status values.
const (
	// Up represents an "up" status.
	Up Status = iota
	// Down represents a "down" status.
	Down
)
