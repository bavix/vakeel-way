package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/bavix/vakeel-way/internal/domain/entities"
)

// StateManager is an interface that defines the behavior for sending status updates
// to a state service.
//
// This interface has a single method:
// - Send: sends a status update for a given UUID.
//
// The Send method has the following parameters:
// - ctx: a context.Context used to cancel the operation if needed.
// - id: a UUID representing the ID of the service.
// - status: a entities.Status representing the status to send.
//
// The Send method returns an error if the status update cannot be sent to the state service,
// and nil if the status update was sent successfully.
//
// The Send method is used to send a status update for a given UUID
// to the state service.
// It takes a context.Context used to cancel the operation if needed,
// a UUID representing the ID of the service,
// and a entities.Status representing the status to send.
// It returns an error if the status update cannot be sent to the state service,
// and nil if the status update was sent successfully.
type StateManager interface {
	// Send sends a status update to the state service for a given UUID.
	//
	// Parameters:
	//   - ctx: The context.Context used to cancel the operation if needed.
	//   - id: The UUID of the service.
	//   - status: The entities.Status to send.
	//
	// Returns:
	//   - An error if the status update cannot be sent to the state service,
	//     or nil if the status update was sent successfully.
	Send(ctx context.Context, id uuid.UUID, status entities.Status) error
}

// Checker represents a struct that handles the logic for sending status updates to the state service.
//
// The Checker struct has the following fields:
// - Events: A channel of type uuid.UUID that is used to send UUIDs to the goroutine that sends status updates.
// - state: A StateManager interface that is used to send status updates to the state service.
type Checker struct {
	// Events is a channel of type uuid.UUID that is used to send UUIDs to the goroutine that sends status updates.
	// The channel has a buffer size of 64.
	Events chan uuid.UUID
	// state is a StateManager interface that is used to send status updates to the state service.
	state StateManager
}

// NewChecker creates a new instance of the Checker struct.
//
// It takes a StateManager interface as a parameter and returns a pointer to a Checker struct.
// The Checker struct is used to handle the logic for sending status updates to the state service.
// It initializes the Events channel with a buffer size of 64, which is used to send UUIDs to
// the goroutine that sends status updates.
//
// Parameters:
//   - client: A StateManager interface used to send events to the state service.
//
// Returns:
//   - A pointer to a Checker struct.
func NewChecker(client StateManager) *Checker {
	const bufferSize = 64 // Buffer size for the Events channel.

	// Create a new instance of the Checker struct.
	// The Checker struct is used to handle the logic for sending status updates to the state service.
	// It initializes the Events channel with a buffer size of 64, which is used to send UUIDs to
	// the goroutine that sends status updates.
	return &Checker{
		// Events is a channel of type uuid.UUID that is used to send UUIDs to the goroutine that sends status updates.
		// The channel has a buffer size of 64.
		Events: make(chan uuid.UUID, bufferSize),
		// state is a StateManager interface that is used to send status updates to the state service.
		state: client,
	}
}

// Send sends an event to the events channel of the Checker.
//
// This function sends an event to the events channel of the Checker,
// which is used to trigger the handler function to process the event.
//
// Parameters:
//   - id: The uuid.UUID object representing the event to be sent.
func (c *Checker) Send(id uuid.UUID) {
	// Send the event to the events channel.
	// The event is sent to the Events channel of the Checker.
	// The Events channel is a channel of type uuid.UUID that is used to send events to the goroutine that processes the events.
	//
	// This function does not return anything.

	// Send the event to the events channel.
	c.Events <- id
}

// Handler is a goroutine that processes events from the Events channel.
//
// This function continuously listens for events on the Events channel.
// When an event is received, it sends a status update to the state service.
// If the context is canceled, the function returns.
//
// Parameters:
// - ctx: The context.Context object that is used to cancel the goroutine.
func (c *Checker) Handler(ctx context.Context) {
	// Get the logger from the context.
	logger := zerolog.Ctx(ctx)

	// Continuously listen for events on the Events channel.
	for {
		// Receive an event from the Events channel.
		// The select statement ensures that the goroutine does not block indefinitely.
		// If the channel is closed, the receive operation will return a boolean value of false.
		select {
		// Receive an event from the Events channel.
		case id, ok := <-c.Events:
			// If the channel is closed, return from the function.
			if !ok {
				return
			}

			// Send a status update to the state service.
			// If an error occurs, log the error.
			if err := c.state.Send(ctx, id, entities.Up); err != nil {
				// Log the error that occurred during sending the event.
				logger.Err(err).Str("id", id.String()).Msg("checker: failed to send event")
			}

		// If the context is canceled, return from the function.
		case <-ctx.Done():
			return
		}
	}
}

// Close closes the Events channel of the Checker.
func (c *Checker) Close() {
	// Close the Events channel to indicate that no more events will be sent.
	close(c.Events)
}
