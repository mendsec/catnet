package output

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mendsec/catnet-core/pkg/engine"
)

type HumanOutput struct {
	noColor bool
	quiet   bool
	writer  *tabwriter.Writer
	out     io.Writer
	errOut  io.Writer
	start   time.Time
}

func NewHumanOutput(noColor, quiet bool) *HumanOutput {
	fi, err := os.Stdout.Stat()
	isTTY := err == nil && (fi.Mode()&os.ModeCharDevice) != 0
	if !isTTY {
		noColor = true
	}

	return newHumanOutputWithWriters(os.Stdout, os.Stderr, noColor, quiet)
}

func newHumanOutputWithWriters(out, errOut io.Writer, noColor, quiet bool) *HumanOutput {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	return &HumanOutput{
		noColor: noColor,
		quiet:   quiet,
		writer:  w,
		out:     out,
		errOut:  errOut,
	}
}

func (h *HumanOutput) HandleEvent(event engine.ScanEvent, total int) {
	switch event.Type {
	case engine.EventLifecycleStart:
		h.start = time.Now()
		if !h.quiet {
			fmt.Fprintf(h.errOut, "Scanning %d hosts...\n\n", total)
			fmt.Fprintln(h.writer, "  IP\tHOSTNAME\tMAC\tSTATUS\tPORTS")
			h.writer.Flush()
		}
	case engine.EventProgress:
		if !h.quiet {
			fmt.Fprintf(h.errOut, "\rProgress: %.1f%%", event.Progress*100)
		}
	case engine.EventResult:
		if event.Device == nil {
			return
		}
		
		status := "DEAD"
		colorStart := ""
		colorEnd := ""
		
		if event.Device.IsAlive {
			status = "ALIVE"
			if !h.noColor {
				colorStart = "\033[32m"
				colorEnd = "\033[0m"
			}
		} else {
			if !h.noColor {
				colorStart = "\033[90m"
				colorEnd = "\033[0m"
			}
		}

		mac := event.Device.MAC
		if mac == "" {
			mac = "—"
		}
		hostname := event.Device.Hostname
		if hostname == "" {
			hostname = "—"
		}

		ports := "—"
		if len(event.Device.OpenPorts) > 0 {
			var pStrs []string
			for _, p := range event.Device.OpenPorts {
				pStrs = append(pStrs, fmt.Sprintf("%d", p))
			}
			ports = strings.Join(pStrs, ", ")
		}

		if !h.quiet {
			fmt.Fprint(h.errOut, "\r                                        \r")
		}

		fmt.Fprintf(h.writer, "  %s\t%s\t%s\t%s%s%s\t%s\n", 
			event.Device.IP, hostname, mac, colorStart, status, colorEnd, ports)
		h.writer.Flush()

	case engine.EventLifecycleComplete:
		if !h.quiet {
			fmt.Fprintf(h.errOut, "\n\nScan complete in %s\n", time.Since(h.start).Round(10*time.Millisecond))
		}
	case engine.EventLifecycleCancel:
		fmt.Fprintf(h.errOut, "\n[CANCELLED] %s\n", event.Message)
	case engine.EventWarning:
		fmt.Fprintf(h.errOut, "\n[WARN] %s\n", event.Message)
	}
}
