package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "vakeel-way",
	Version: version,
	Short:   "A brief description of your application",
	RunE: func(cmd *cobra.Command, args []string) error {
		toggle, err := cmd.Flags().GetBool("toggle")
		if err != nil {
			return err
		}

		if toggle {
			cmd.Printf("hello world\n")
		}

		return nil
	},
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
