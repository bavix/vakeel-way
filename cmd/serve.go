package cmd

import (
	"errors"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/bavix/vakeel-way/internal/build"
	"github.com/bavix/vakeel-way/internal/config"
)

var cfgFile string

// serveCmd returns the serve command.
//
// The serve command starts the server. The server is a gRPC service that
// provides the StateService RPC service.
//
//nolint:exhaustruct
func serveCmd() *cobra.Command {
	// Create a new serve command.
	return &cobra.Command{
		Use:   "serve",
		Short: "Starts the server",
		// RunE is the function that is called when the command is executed.
		// It returns an error if there is a problem starting the server.
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Create a new context that listens for the interrupt signal.
			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer cancel()

			// Read the configuration from the environment variables.
			cfg, err := config.New(cfgFile)
			if err != nil {
				return err
			}

			// Create a new builder using the configuration.
			builder, err := build.NewBuilder(cfg)
			if err != nil {
				return err
			}

			// Run the gRPC server using the builder. The context is used to log
			// messages related to the gRPC server.
			if err := builder.RunGRPCServer(
				builder.Logger(ctx),
			); !errors.Is(err, grpc.ErrServerStopped) {
				return err
			}

			// Return nil if the server is stopped successfully.
			return nil
		},
	}
}

// init is a special Go function that is called after all the variable
// declarations in the package have evaluated their initializers.
//
// In this function, we are adding the serve command to the root command.
// The serve command starts the server.
// The server is a gRPC service that provides the StateService
// RPC service.
//
// This function is called automatically by the Go runtime.
func init() {
	// Create the serve command.
	serveCmd := serveCmd()

	// Add the serve command to the root command.
	rootCmd.AddCommand(serveCmd)

	// Add a flag to the root command that specifies the location of the
	// configuration file.
	serveCmd.Flags().StringVar(
		&cfgFile,
		"config",
		"/etc/vakeel-way/config.yaml",
		"Path to the configuration file.",
	)
}
