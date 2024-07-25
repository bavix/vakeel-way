package build

import (
	"context"
	"log"
	"time"

	"github.com/rs/zerolog"
)

// Logger creates a new context with a logger attached to it.
//
// It creates a logger with the log level specified in the configuration file.
// The logger is then attached to the given context.
//
// Parameters:
//   - ctx: The context to attach the logger to.
//
// Returns:
//   - The context with the logger attached.
func (b *Builder) Logger(ctx context.Context) context.Context {
	// Parse the log level from the configuration file.
	level, err := zerolog.ParseLevel(b.config.Log.Level)
	if err != nil {
		// If the log level is invalid, log the error and stop the application.
		log.Fatal(err)
	}

	// Create a new logger with the specified log level and time format.
	// The time format is set to RFC3339Nano, which is the most precise time format.
	logger := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339Nano
	})).
		Level(level).
		With().
		Timestamp().
		Logger()

	// Attach the logger to the given context and return the new context.
	return logger.WithContext(ctx)
}
