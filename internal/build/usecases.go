package build

import (
	"context"

	"github.com/bavix/vakeel-way/internal/domain/services"
	"github.com/bavix/vakeel-way/internal/domain/usecases"
)

// checkerUsecase returns a new instance of the Checker usecase.
// If the Builder instance already has a Checker instance, it will be returned.
// Otherwise, a new Checker instance will be created and stored in the Builder instance.
//
// Parameters:
//   - ctx: The context.Context used to cancel the operation if needed.
//
// Returns:
//   - A pointer to a Checker usecase.
func (b *Builder) checkerUsecase(ctx context.Context) *usecases.Checker {
	// Check if the Builder instance already has a Checker instance.
	if b.checker != nil {
		return b.checker
	}

	// Create a new StateManager instance.
	// The StateManager instance is responsible for sending status updates to the state service.
	// It takes a context.Context used to cancel the operation if needed,
	// a WebhookRepository instance used to retrieve webhooks by their UUIDs,
	// and an InStatusClient instance that is used to send status updates to the state service.
	stateManager := services.NewStateManager(
		b.inStatusClient(),
		b.WebhookRepository(),
	)

	// Create a new Checker instance using the StateManager instance.
	// The Checker instance is responsible for sending status updates to the state service.
	// It takes a StateManager instance as a parameter.
	b.checker = usecases.NewChecker(stateManager)

	// Start a goroutine to close the Checker instance when the context is canceled.
	// This ensures that the Checker goroutine is stopped when the context is canceled.
	go func() {
		<-ctx.Done()
		b.checker.Close()
	}()

	// Start a goroutine to process events from the Checker's Events channel.
	// The goroutine listens for events on the Events channel and sends status updates to the state service.
	// If the context is canceled, the goroutine returns.
	go b.checker.Handler(ctx)

	return b.checker
}
