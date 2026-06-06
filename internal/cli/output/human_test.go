package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mendsec/catnet-core/pkg/engine"
	"github.com/mendsec/catnet-core/pkg/results"
)

func TestHumanOutputFormat(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	h := newHumanOutputWithWriters(out, errOut, true, false)

	h.HandleEvent(engine.ScanEvent{
		Type: engine.EventLifecycleStart,
	}, 1)

	if !strings.Contains(errOut.String(), "Scanning 1 hosts") {
		t.Errorf("Expected 'Scanning 1 hosts' in stderr, got: %s", errOut.String())
	}
	if !strings.Contains(out.String(), "IP") || !strings.Contains(out.String(), "HOSTNAME") {
		t.Errorf("Expected table header in stdout, got: %s", out.String())
	}

	h.HandleEvent(engine.ScanEvent{
		Type: engine.EventResult,
		Device: &results.DeviceInfo{
			IP:        "192.168.1.1",
			Hostname:  "router",
			MAC:       "AA:BB:CC:DD:EE:FF",
			IsAlive:   true,
			OpenPorts: []int{80, 443},
		},
	}, 1)

	outStr := out.String()
	if !strings.Contains(outStr, "192.168.1.1") ||
		!strings.Contains(outStr, "router") ||
		!strings.Contains(outStr, "AA:BB:CC:DD:EE:FF") ||
		!strings.Contains(outStr, "ALIVE") ||
		!strings.Contains(outStr, "80, 443") {
		t.Errorf("Expected output to contain device details, got: %s", outStr)
	}
}

func TestHumanOutputQuiet(t *testing.T) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	h := newHumanOutputWithWriters(out, errOut, true, true)

	h.HandleEvent(engine.ScanEvent{
		Type: engine.EventLifecycleStart,
	}, 1)

	if len(errOut.String()) > 0 {
		t.Errorf("Expected empty errOut in quiet mode, got: %s", errOut.String())
	}
}
