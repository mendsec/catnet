package output

import (
	"fmt"
	"io"
	"os"

	"github.com/mendsec/catnet-core/pkg/engine"
)

type JSONOutput struct {
	quiet  bool
	errOut io.Writer
}

func NewJSONOutput(quiet bool) *JSONOutput {
	return &JSONOutput{quiet: quiet, errOut: os.Stderr}
}

func newJSONOutputWithWriter(errOut io.Writer, quiet bool) *JSONOutput {
	return &JSONOutput{quiet: quiet, errOut: errOut}
}

func (j *JSONOutput) HandleEvent(event engine.ScanEvent, total int) {
	switch event.Type {
	case engine.EventLifecycleStart:
		if !j.quiet {
			fmt.Fprintf(j.errOut, "Scanning %d hosts...\n", total)
		}
	case engine.EventProgress:
		if !j.quiet {
			fmt.Fprintf(j.errOut, "\rProgress: %.1f%%", event.Progress*100)
		}
	case engine.EventLifecycleComplete:
		if !j.quiet {
			fmt.Fprintln(j.errOut, "\nScan complete.")
		}
	case engine.EventLifecycleCancel:
		fmt.Fprintf(j.errOut, "\n[CANCELLED] %s\n", event.Message)
	case engine.EventWarning:
		fmt.Fprintf(j.errOut, "\n[WARN] %s\n", event.Message)
	}
}
