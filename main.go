package main

import (
	"context"

	"github.com/bavix/vakeel-way/cmd"
)

// main is the entry point of the application.
//
// It sets up the command-line flags and executes the root command with the
// given context.
func main() {
	// Create a new context with no parent.
	// This context is used as the parent for all child contexts.
	ctx := context.Background()

	// Execute the root command with the created context.
	//
	// The root command is the main command of the application.
	// It sets up the command-line flags and starts the application.
	// If there is an error during the execution of the command,
	// it will exit the program with a status code of 1.
	cmd.Execute(ctx)
}
