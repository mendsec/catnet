package cli

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var Version = "dev"
var shortVersion bool

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		if shortVersion {
			fmt.Println(Version)
			return
		}

		fmt.Printf("catnet/%s (%s/%s) %s\n", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())

		coreVersion := "unknown"
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			for _, dep := range buildInfo.Deps {
				if dep.Path == "github.com/mendsec/catnet-core" {
					coreVersion = dep.Version
					break
				}
			}
		}
		fmt.Printf("catnet-core/%s\n", coreVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&shortVersion, "short", false, "Print only the version number")
}
