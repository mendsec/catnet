package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mendsec/catnet-core/pkg/engine"
	"github.com/mendsec/catnet-core/pkg/exporter"
	"github.com/mendsec/catnet-core/pkg/targets"
	"github.com/mendsec/catnet/internal/cli/output"
)

var scanCmd = &cobra.Command{
	Use:   "scan [targets]",
	Short: "Scan a network range for live hosts",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		scanQuiet, _ := cmd.Flags().GetBool("quiet")
		noColor, _ := cmd.Flags().GetBool("no-color")
		scanOutput, _ := cmd.Flags().GetString("output")
		scanNoPorts, _ := cmd.Flags().GetBool("no-ports")

		var allIPs []string
		for _, arg := range args {
			parts := strings.Split(arg, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				ips, err := targets.ParseRange(part)
				if err != nil {
					return NewExitError(ExitCodeInputError, "Invalid target '%s': %v", part, err)
				}
				allIPs = append(allIPs, ips...)
			}
		}

		if len(allIPs) == 0 {
			return NewExitError(ExitCodeInputError, "No valid targets provided.")
		}

		cfg := engine.DefaultConfig()

		if scanNoPorts {
			cfg.DefaultPorts = []int{}
		} else if cmd.Flags().Changed("ports") {
			cfg.DefaultPorts, _ = cmd.Flags().GetIntSlice("ports")
		}

		if cmd.Flags().Changed("threads") {
			cfg.MaxThreads, _ = cmd.Flags().GetInt("threads")
		}
		if cmd.Flags().Changed("ping-timeout") {
			cfg.PingTimeoutMs, _ = cmd.Flags().GetInt("ping-timeout")
		}
		if cmd.Flags().Changed("port-timeout") {
			cfg.PortTimeoutMs, _ = cmd.Flags().GetInt("port-timeout")
		}
		cfg.Sanitize()

		ctx, cancel := WithCancelOnSignal(cmd.Context())
		defer cancel()

		var eventHandler engine.EventCallback

		switch format {
		case "json":
			handler := output.NewJSONOutput(scanQuiet)
			eventHandler = func(e engine.ScanEvent) {
				handler.HandleEvent(e, len(allIPs))
			}
		case "human":
			handler := output.NewHumanOutput(noColor, scanQuiet)
			eventHandler = func(e engine.ScanEvent) {
				handler.HandleEvent(e, len(allIPs))
			}
		default:
			return NewExitError(ExitCodeInputError, "Unsupported format '%s'. Use 'json' or 'human'.", format)
		}

		report, err := engine.StartScan(ctx, allIPs, cfg, eventHandler)

		if err != nil {
			if ctx.Err() != nil {
				return NewExitError(ExitCodeInterrupted, "Scan cancelled")
			}
			return NewExitError(ExitCodeRuntimeError, "Scan failed: %v", err)
		}

		if format == "json" {
			jsonBytes, err := exporter.ExportJSON(report)
			if err != nil {
				return NewExitError(ExitCodeRuntimeError, "Failed to encode JSON: %v", err)
			}

			if scanOutput != "" {
				if err := os.WriteFile(scanOutput, jsonBytes, 0600); err != nil {
					return NewExitError(ExitCodeRuntimeError, "Failed to write output file: %v", err)
				}
			} else {
				fmt.Print(string(jsonBytes))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443, 139, 445, 3389}, "Ports to scan")
	scanCmd.Flags().IntP("threads", "t", 64, "Max concurrent threads")
	scanCmd.Flags().Int("ping-timeout", 1000, "Ping timeout in milliseconds")
	scanCmd.Flags().Int("port-timeout", 500, "Port timeout in milliseconds")
	scanCmd.Flags().StringP("output", "o", "", "Write JSON output to file instead of stdout")
	scanCmd.Flags().BoolP("quiet", "q", false, "Suppress progress output")
	scanCmd.Flags().Bool("no-ports", false, "Skip port scanning entirely")
	scanCmd.Flags().String("format", "human", "Output format: json, human")
}
