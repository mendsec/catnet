package catnet

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mendsec/catnet-core/pkg/events"
	"github.com/mendsec/catnet-core/pkg/export"
	"github.com/mendsec/catnet-core/pkg/profile"
	"github.com/mendsec/catnet-core/pkg/results"
	"github.com/mendsec/catnet-core/pkg/scan"
	"github.com/mendsec/catnet-core/pkg/targets"
	"github.com/spf13/cobra"
)

var (
	exportFormat string
	threads      int
	timeoutMs    int
	verbose      bool
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&exportFormat, "export", "e", "text", "Export format (json, csv, text)")
	scanCmd.Flags().IntVarP(&threads, "threads", "t", 64, "Max concurrent threads")
	scanCmd.Flags().IntVar(&timeoutMs, "timeout", 1000, "Ping and port timeout in ms")
	scanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

var scanCmd = &cobra.Command{
	Use:   "scan [targets]",
	Short: "Run a network scan against specified targets",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetInput := args[0]
		ips, err := targets.ParseRange(targetInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse targets: %v\n", err)
			os.Exit(1)
		}

		engine := scan.NewEngine()
		prof := profile.DefaultProfile()
		prof.Concurrency = threads
		prof.TimeoutMs = timeoutMs

		eventChan := make(chan events.Event)
		done := make(chan struct{})
		var scannedHosts []results.HostResult

		exportFormat = strings.ToLower(exportFormat)
		isTextual := exportFormat == "text"

		go func() {
			defer close(done)
			for ev := range eventChan {
				switch ev.Type {
				case events.ScanStarted:
					if isTextual || verbose {
						fmt.Fprintf(os.Stderr, "[*] Scan started on %d targets (Threads: %d, Timeout: %dms)\n", ev.Data, threads, timeoutMs)
					}
				case events.HostDiscovered:
					data, ok := ev.Data.(events.HostDiscoveredData)
					if ok {
						scannedHosts = append(scannedHosts, data.Host)
						if verbose {
							if data.Host.Alive {
								fmt.Fprintf(os.Stderr, "[+] Host UP: %s (MAC: %s) - Ports: %v\n", data.Host.IP, data.Host.MAC, data.Host.OpenPorts)
							} else {
								fmt.Fprintf(os.Stderr, "[-] Host DOWN: %s\n", data.Host.IP)
							}
						}
					}
				case events.ScanProgress:
					if verbose {
						data, ok := ev.Data.(events.ProgressData)
						if ok && data.Processed%10 == 0 {
							fmt.Fprintf(os.Stderr, "[~] Progress: %d/%d (%.1f%%)\n", data.Processed, data.Total, data.Ratio*100)
						}
					}
				case events.ScanCompleted:
					if isTextual || verbose {
						fmt.Fprintln(os.Stderr, "[*] Scan completed")
					}
				}
			}
		}()

		err = engine.ScanStream(context.Background(), ips, prof, eventChan)
		close(eventChan)
		<-done // Wait for the event loop to finish printing

		if err != nil {
			fmt.Fprintf(os.Stderr, "Scan failed: %v\n", err)
			os.Exit(1)
		}

		if exportFormat == "json" {
			jsonBytes, _ := export.ExportJSON(scannedHosts)
			fmt.Println(string(jsonBytes))
			return
		}

		if exportFormat == "csv" {
			csvBytes, _ := export.ExportCSV(scannedHosts)
			fmt.Print(string(csvBytes))
			return
		}

		// Text format (default)
		fmt.Println("\n--- Scan Results ---")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "IP\tHOSTNAME\tMAC\tSTATUS\tOPEN PORTS")
		for _, h := range scannedHosts {
			status := "Dead"
			if h.Alive {
				status = "Alive"
			}
			ports := fmt.Sprintf("%v", h.OpenPorts)
			if len(h.OpenPorts) == 0 {
				ports = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", h.IP, h.Hostname, h.MAC, status, ports)
		}
		w.Flush()
	},
}
