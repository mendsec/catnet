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

var (
	scanPorts       []int
	scanThreads     int
	scanPingTimeout int
	scanPortTimeout int
	scanOutput      string
	scanQuiet       bool
	scanNoPorts     bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [targets]",
	Short: "Scan a network range for live hosts",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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
			cfg.DefaultPorts = scanPorts
		}
		
		cfg.MaxThreads = scanThreads
		cfg.PingTimeoutMs = scanPingTimeout
		cfg.PortTimeoutMs = scanPortTimeout
		cfg.Sanitize()

		ctx, cancel := WithCancelOnSignal(cmd.Context())
		defer cancel()

		var eventHandler engine.EventCallback
		
		if format == "json" {
			handler := output.NewJSONOutput(scanQuiet)
			eventHandler = func(e engine.ScanEvent) {
				handler.HandleEvent(e, len(allIPs))
			}
		} else {
			handler := output.NewHumanOutput(noColor, scanQuiet)
			eventHandler = func(e engine.ScanEvent) {
				handler.HandleEvent(e, len(allIPs))
			}
		}

		report, err := engine.StartScan(ctx, allIPs, cfg, eventHandler)
		
		if err != nil {
			if ctx.Err() != nil {
				return NewExitError(ExitCodeInterrupted, "Scan cancelled")
			}
			return NewExitError(ExitCodeRuntimeError, "Scan failed: %v", err)
		}

		if format == "json" || scanOutput != "" {
			jsonBytes, err := exporter.ExportJSON(report)
			if err != nil {
				return NewExitError(ExitCodeRuntimeError, "Failed to encode JSON: %v", err)
			}
			
			if scanOutput != "" {
				if err := os.WriteFile(scanOutput, jsonBytes, 0644); err != nil {
					return NewExitError(ExitCodeRuntimeError, "Failed to write output file: %v", err)
				}
			} else {
				fmt.Println(string(jsonBytes))
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceVarP(&scanPorts, "ports", "p", []int{22, 80, 443, 139, 445, 3389}, "Ports to scan")
	scanCmd.Flags().IntVarP(&scanThreads, "threads", "t", 64, "Max concurrent threads")
	scanCmd.Flags().IntVar(&scanPingTimeout, "ping-timeout", 1000, "Ping timeout in milliseconds")
	scanCmd.Flags().IntVar(&scanPortTimeout, "port-timeout", 500, "Port timeout in milliseconds")
	scanCmd.Flags().StringVarP(&scanOutput, "output", "o", "", "Write JSON output to file instead of stdout")
	scanCmd.Flags().BoolVarP(&scanQuiet, "quiet", "q", false, "Suppress progress output")
	scanCmd.Flags().BoolVar(&scanNoPorts, "no-ports", false, "Skip port scanning entirely")
}
