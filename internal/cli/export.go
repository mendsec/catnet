package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mendsec/catnet-core/pkg/exporter"
	"github.com/mendsec/catnet-core/pkg/results"
)

var (
	exportOutput string
)

var exportCmd = &cobra.Command{
	Use:   "export [input.json]",
	Short: "Export a previous scan result to a different format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		data, err := os.ReadFile(inputFile)
		if err != nil {
			return NewExitError(ExitCodeInputError, "Failed to read input file: %v", err)
		}

		var report results.ScanReport
		// Usar json.Unmarshal diretamente permite que campos desconhecidos sejam ignorados (DisallowUnknownFields é falso por default)
		if err := json.Unmarshal(data, &report); err != nil {
			return NewExitError(ExitCodeInputError, "Failed to parse JSON: %v", err)
		}

		if report.SchemaVersion != "" {
			if len(report.SchemaVersion) > 0 && report.SchemaVersion[0] != '1' {
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
			if err := os.WriteFile(exportOutput, outBytes, 0644); err != nil {
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
	// We override the persistent 'format' flag for the export command, adding a shorthand and different default.
	exportCmd.Flags().StringP("format", "f", "", "Output format: json, csv, xml (required)")
	exportCmd.MarkFlagRequired("format")
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Write output to file (default: stdout)")
}
