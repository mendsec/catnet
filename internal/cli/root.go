package cli

import (
	"github.com/spf13/cobra"
)

var (
	format  string
	noColor bool
)

var rootCmd = &cobra.Command{
	Use:   "catnet",
	Short: "catnet — Network scanner CLI",
	Long:  "catnet — Network scanner CLI",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&format, "format", "human", "Output format: json, human")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
