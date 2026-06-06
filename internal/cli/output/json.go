package output

import (
	"fmt"
	"os"

	"github.com/mendsec/catnet-core/pkg/engine"
)

type JSONOutput struct {
	quiet bool
}

func NewJSONOutput(quiet bool) *JSONOutput {
	return &JSONOutput{quiet: quiet}
}

func (j *JSONOutput) HandleEvent(event engine.ScanEvent, total int) {
	switch event.Type {
	case engine.EventLifecycleStart:
		if !j.quiet {
			fmt.Fprintf(os.Stderr, "Scanning %d hosts...\n", total)
		}
	case engine.EventProgress:
		if !j.quiet {
			fmt.Fprintf(os.Stderr, "\rProgress: %.1f%%", event.Progress*100)
		}
	case engine.EventLifecycleComplete:
		if !j.quiet {
			fmt.Fprintln(os.Stderr, "\nScan complete.")
		}
	case engine.EventLifecycleCancel:
		fmt.Fprintf(os.Stderr, "\n[CANCELLED] %s\n", event.Message)
	case engine.EventWarning:
		fmt.Fprintf(os.Stderr, "\n[WARN] %s\n", event.Message)
	}
}
