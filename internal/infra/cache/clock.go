package cache

import "time"

// Click is an interface that provides the current time.
//
// It defines the Now method that returns the current local time.
type Click interface {
	// Now returns the current local time.
	Now() time.Time
}

// clock is a struct that implements the Clock interface.
//
// It provides the current time.
type clock struct{}

// Now is a method that implements the Clock interface.
//
// It returns the current local time.
//
// Returns:
//
//	time.Time: The current local time.
func (c clock) Now() time.Time {
	// Return the current time.
	return time.Now()
}
