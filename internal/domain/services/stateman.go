package services

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/bavix/vakeel-way/internal/domain/entities"
	"github.com/bavix/vakeel-way/internal/infra/cache"
)

// WebhookRegistry represents an interface for managing webhooks.
//
// A WebhookRegistry is responsible for retrieving webhooks by their IDs.
// It provides a Get method for retrieving a webhook by its UUID.
// The Get method takes a context.Context used to cancel the operation if needed
// and a UUID representing the ID of the webhook to retrieve.
// It returns a string representing the webhook data and an error if the webhook
// is not found or if there is an issue retrieving it.
type WebhookRegistry interface {
	// Get retrieves a webhook by its ID.
	//
	// Parameters:
	//   - ctx: The context.Context used to cancel the operation if needed.
	//   - id: The UUID of the webhook to retrieve.
	//
	// Returns:
	//   - webhookData: The webhook data.
	//   - err: An error if the webhook is not found or if there is an issue retrieving it.
	Get(ctx context.Context, id uuid.UUID) (webhookData string, err error)

	// All returns all webhook IDs.
	//
	// This method returns all webhook IDs as a slice of UUIDs.
	All() []uuid.UUID
}

// Api represents an interface for sending status updates.

// API represents an interface for sending status updates.
type API interface {
	// Send sends a status update to the specified URL.
	//
	// Parameters:
	//   - ctx: The context.Context used to cancel the operation if needed.
	//   - url: The URL to send the status update to.
	//   - status: The entities.Status to send.
	//
	// Returns:
	//   - An error if the status update cannot be sent to the URL.
	//   - nil if the status update was sent successfully.
	//
	// Send sends a status update to the specified URL.
	// It takes a context.Context used to cancel the operation if needed,
	// a string representing the URL to send the status update to,
	// and an entities.Status representing the status to send.
	// It returns an error if the status update cannot be sent to the URL,
	// and nil if the status update was sent successfully.
	Send(ctx context.Context, url string, status entities.Status) error
}

// state represents the current status of a webhook.
//
// The state struct holds the current status of a webhook. It has the following fields:
//   - status: The current status of the webhook.
//   - attempt: The number of attempts made to send a status update to the webhook.
type state struct {
	// status is the current status of the webhook.
	status entities.Status

	// attempt is the number of attempts made to send a status update to the webhook.
	attempt uint32
}

// StateManager manages the sending of status updates to webhooks.
//
// The StateManager struct holds the necessary dependencies to manage the sending of status updates to webhooks.
// It has the following fields:
//   - api: The API used to send status updates.
//   - repo: The repository used to get webhook URLs.
//   - cache: The cache used to store the current status of webhooks.
//   - mu: The mutex used to synchronize access to the cache.
type StateManager struct {
	// api is the API used to send status updates.
	//
	// This field holds the API used to send status updates. It is of type Api.
	api API

	// repo is the repository used to get webhook URLs.
	//
	// This field holds the repository used to get webhook URLs. It is of type WebhookRegistry.
	repo WebhookRegistry

	// cache is the cache used to store the current status of webhooks.
	//
	// This field holds the cache used to store the current status of webhooks.
	// It is of type *cache.Cache[uuid.UUID, state].
	cache *cache.Cache[uuid.UUID, state]

	// mu is the mutex used to synchronize access to the cache.
	//
	// This field holds the mutex used to synchronize access to the cache.
	// It is of type sync.RWMutex.
	mu sync.RWMutex

	// log is the logger used to log messages related to the StateManager.
	//
	// This field holds the logger used to log messages related to the StateManager.
	// It is of type *zerolog.Logger.
	log *zerolog.Logger
}

// NewStateManager creates a new instance of the StateManager struct.
//
// It takes an API, a WebhookRegistry, and a logger as input parameters.
// It returns a pointer to the initialized StateManager.
//
// Parameters:
//   - api: The API used to send status updates.
//   - repo: The repository used to get webhook URLs.
//   - log: The logger used to log messages.
//
// Returns:
//   - A pointer to the initialized StateManager.
//
//nolint:exhaustruct
func NewStateManager(api API, repo WebhookRegistry, log *zerolog.Logger) *StateManager {
	// Create a new StateManager instance.
	stateManager := &StateManager{
		api:  api,  // Set the API used to send status updates.
		repo: repo, // Set the repository used to get webhook URLs.
		log:  log,  // Set the logger used to log messages.
	}

	// Create a new cache with a length based on the number of webhooks.
	// The cache is initialized with the garbage collector function set to
	// garbageCollector.
	cache := cache.NewCache(
		len(repo.All()), // Initialize the cache size.
		cache.WithOnEvict(stateManager.garbageCollector), // Set the garbage collector function.
	)

	// Assign the cache to the StateManager instance.
	stateManager.cache = cache

	// Return the initialized StateManager.
	return stateManager
}

// garbageCollector is a function that is called when an item is evicted from the cache.
// It sends a status update to the specified webhook URL if the status is different
// from the current status in the cache.
//
// Parameters:
//   - id: The UUID of the webhook.
//   - current: The current state of the webhook in the cache.
func (s *StateManager) garbageCollector(id uuid.UUID, current state) {
	// Maximum number of attempts to send a status update.
	const maxAttempts = 5

	// Check if the maximum number of attempts has been reached.
	if current.attempt >= maxAttempts {
		return
	}

	// Lock the mutex to ensure exclusive access to the cache.
	s.mu.Lock()
	defer s.mu.Unlock()

	// Set a timeout for the operation.
	const timeout = 15 * time.Second

	// Create a context with the timeout.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Get the URL of the webhook from the repository.
	target, err := s.repo.Get(ctx, id)
	if err != nil {
		// Increment the number of attempts.
		atomic.AddUint32(&current.attempt, 1)

		// If an error occurs, add the status 'Down' to the cache.
		s.cache.Add(id, current, timeout)

		return
	}

	// Inform the webhook about the status update.
	s.inform(id, entities.Down)

	// Send a status update to the URL.
	err = s.api.Send(ctx, target, entities.Down)
	if err != nil {
		// Increment the number of attempts.
		atomic.AddUint32(&current.attempt, 1)

		// If an error occurs, add the status 'Down' to the cache.
		s.cache.Add(id, current, timeout)

		return
	}
}

// Send sends a status update to the specified webhook ID.
//
// If the status is the same as the current status in the cache,
// the status update is not sent and the status is added to the cache.
//
// Parameters:
//   - ctx: The context.Context used to cancel the operation if needed.
//   - id: The UUID of the webhook.
//   - status: The entities.Status to send.
//
// Returns:
//   - An error if the webhook URL cannot be retrieved from the repository,
//     or if the status update cannot be sent to the webhook.
//   - nil if the status update was sent successfully or if the status is the
//     same as the current status in the cache.
func (s *StateManager) Send(ctx context.Context, id uuid.UUID, status entities.Status) error {
	// The TTL (Time to Live) of the status in the cache.
	const ttl = time.Minute

	// Get the current status from the cache.
	currentStatus, _ := s.cache.Get(id)

	// If the status is the same as the current status in the cache,
	// add it to the cache and return nil.
	if currentStatus != nil && currentStatus.status == status {
		// Prolong the life of the status in the cache.
		s.cache.Add(id, state{status: status, attempt: 0}, ttl)

		return nil
	}

	// Get the webhook URL from the repository.
	// This is the URL of the webhook that will receive the status update.
	target, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	// Inform the logger that a status update is being sent.
	// This logs the ID and status of the service being updated.
	s.inform(id, status)

	// Send the status update to the webhook.
	// This sends a POST request to the webhook URL with the status as the request body.
	if err := s.api.Send(ctx, target, status); err != nil {
		return err
	}

	// Add the status to the cache.
	// This adds the status to the cache so that it can be retrieved later.
	s.cache.Add(id, state{status: status, attempt: 0}, ttl)

	return nil
}

// inform logs the sending of a status update.
//
// It logs the ID and status of the service being updated.
// It takes the ID of the service and its status as parameters.
func (s *StateManager) inform(id uuid.UUID, status entities.Status) {
	// Log the sending of a status update.
	//
	// The log message includes the ID and status of the service being updated.
	// It takes the ID of the service and its status as parameters.
	s.log.Info().
		// The ID of the service.
		Str("id", id.String()).
		// The status of the service.
		Str("status", status.String()).
		// The message to log.
		Msg("Sending status update")
}
