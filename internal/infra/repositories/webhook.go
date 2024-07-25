package repositories

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

// ErrWebhookNotFound is an error that indicates that the requested webhook was not found.
var ErrWebhookNotFound = errors.New("webhook not found")

// WebhookStubRepository is a simple in-memory implementation of the WebhookRepository interface.
//
// It stores the UUIDs and their associated values in a map. The mutex is used to synchronize access to the map.
//
// Fields:
//
// storage is a map that stores the UUIDs and their associated values.
// The map is used to store the UUIDs as keys and their associated values as values.
//
// mu is a mutex used to synchronize access to the storage map.
// The mutex is used to ensure that only one goroutine can modify the storage map at a time.
type WebhookStubRepository struct {
	// storage is a map that stores the UUIDs and their associated values.
	storage map[uuid.UUID]string
	// mu is a mutex used to synchronize access to the storage map.
	// The mutex is used to ensure that only one goroutine can modify the storage map at a time.
	mu sync.Mutex
}

// NewWebhookRepository creates a new instance of the WebhookStubRepository.
//
// This function takes a map that stores the UUIDs and their associated values as input and returns
// a pointer to the newly created WebhookStubRepository.
// The WebhookStubRepository is a simple in-memory implementation of the WebhookRepository interface,
// which stores the UUIDs and their associated values in a map.
//
// Parameters:
// - storage: A map that stores the UUIDs and their associated values.
//
// Returns:
// - A pointer to the newly created WebhookStubRepository.
//
//nolint:exhaustruct
func NewWebhookRepository(storage map[uuid.UUID]string) *WebhookStubRepository {
	// Create a new instance of the WebhookStubRepository.
	// The WebhookStubRepository stores the UUIDs and their associated values in the provided map.
	return &WebhookStubRepository{
		storage: storage, // Store the UUIDs and their associated values in the storage map.
	}
}

// Get retrieves the value associated with the given UUID from the storage.
//
// Parameters:
// - ctx: The context.Context used to cancel the operation if needed.
// - id: The UUID of the webhook.
//
// Returns:
// - The value associated with the given UUID.
// - An error if the UUID is not found.
//
// The function locks the mutex to prevent concurrent access to the storage.
// It retrieves the value associated with the given UUID from the storage.
// If the UUID is not found, it returns an error.
// Otherwise, it returns the value associated with the given UUID.
func (w *WebhookStubRepository) Get(_ context.Context, id uuid.UUID) (string, error) {
	// Lock the mutex to prevent concurrent access to the storage.
	w.mu.Lock()
	defer w.mu.Unlock()

	// Retrieve the value associated with the given UUID from the storage.
	val, ok := w.storage[id]

	// If the UUID is not found, return an error.
	if !ok {
		// Return an error indicating that the webhook was not found.
		return "", ErrWebhookNotFound
	}

	// Return the value associated with the given UUID.
	return val, nil
}

// All returns all keys from the storage.
//
// This function returns all keys from the storage as a slice of UUIDs.
// It locks the mutex to prevent concurrent access to the storage.
//
// Parameters:
// - None
//
// Returns:
// - A slice of UUIDs.
func (w *WebhookStubRepository) All() []uuid.UUID {
	// Lock the mutex to prevent concurrent access to the storage.
	// This is done to prevent data races between goroutines.
	w.mu.Lock()
	defer w.mu.Unlock()

	// Create a slice to store the keys.
	// The length of the slice is set to the capacity of the storage
	// to avoid resizing during iteration.
	// This is done to optimize the performance of the function.
	keys := make([]uuid.UUID, 0, len(w.storage))

	// Iterate over the storage and append each key to the slice.
	// This is done to retrieve all the keys from the storage.
	for key := range w.storage {
		// Append the key to the slice.
		// This is done to build the slice of keys.
		keys = append(keys, key)
	}

	// Return the slice of keys.
	// This is done to return the result of the function.
	return keys
}
