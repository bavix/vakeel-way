package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

// version is the version of the application.
//
// It is used in the 'vakeel-way --version' command and as part of the
// application's version.
var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vakeel-way",               // The name of the command
	Version: version,                    // The version of the command
	Short:   "Collector storage server", // A brief description of the command
	// Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:
	//
	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
}

// Execute runs the root command with the given context.
//
// It executes the root command, which is the main entry point of the application.
// It takes a context.Context as a parameter which is used to cancel or control the execution of the function.
// If there is an error during the execution of the command, the function will exit the program with a status code of 1.
func Execute(ctx context.Context) {
	// Execute the root command with the given context.
	// If there is an error, it will exit the program with a status code of 1.
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
