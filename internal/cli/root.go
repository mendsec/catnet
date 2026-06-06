package cli

import (
	"github.com/spf13/cobra"
)

var (
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
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
