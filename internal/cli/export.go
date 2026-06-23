package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mendsec/catnet-core/pkg/exporter"
	"github.com/mendsec/catnet-core/pkg/results"
)

var exportCmd = &cobra.Command{
	Use:   "export [input.json]",
	Short: "Export a previous scan result to a different format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		exportOutput, _ := cmd.Flags().GetString("output")
		inputFile := args[0]
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return NewExitError(ExitCodeInputError, "Failed to read input file: %v", err)
		}

		var report results.ScanReport
		if err := json.Unmarshal(data, &report); err != nil {
			return NewExitError(ExitCodeInputError, "Failed to parse JSON: %v", err)
		}

		if report.SchemaVersion != "" {
			majorStr, _, _ := strings.Cut(report.SchemaVersion, ".")
			major, err := strconv.Atoi(majorStr)
			if err != nil || major < 1 || major > 2 {
				fmt.Fprintf(os.Stderr, "[WARN] Unknown schema version '%s'. Export might skip unknown fields.\n", report.SchemaVersion)
			}
		}

		fVal, _ := cmd.Flags().GetString("format")

		var outBytes []byte
		switch fVal {
		case "json":
			outBytes, err = exporter.ExportJSON(&report)
		case "csv":
			outBytes, err = exporter.ExportCSV(&report)
		case "xml":
			outBytes, err = exporter.ExportXML(&report)
		default:
			return NewExitError(ExitCodeInputError, "Unsupported export format: %s", fVal)
		}

		if err != nil {
			return NewExitError(ExitCodeRuntimeError, "Failed to export %s: %v", fVal, err)
		}

		if exportOutput != "" {
			if err := os.WriteFile(exportOutput, outBytes, 0600); err != nil {
				return NewExitError(ExitCodeRuntimeError, "Failed to write output file: %v", err)
			}
		} else {
			fmt.Println(string(outBytes))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("format", "f", "", "Output format: json, csv, xml (required)")
	exportCmd.MarkFlagRequired("format")
	exportCmd.Flags().StringP("output", "o", "", "Write output to file (default: stdout)")
}
