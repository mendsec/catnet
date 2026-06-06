package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/mendsec/catnet/internal/cli"
)

func resetFlags() {
	// Root flag reset
	os.Args = []string{"catnet"}
	// It's a bit tricky to reset Cobra flags cleanly between tests if they run in same process
}

func TestScanOutputJSON(t *testing.T) {
	os.Args = []string{"catnet", "scan", "127.0.0.1", "--format", "json", "--quiet"}
	
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cli.Execute()
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
	
	if ver, ok := data["schemaVersion"].(string); !ok || ver != "1.0.0" {
		t.Errorf("Expected schemaVersion 1.0.0, got %v", ver)
	}
}

func TestScanCancelledByContext(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Signal testing is unreliable on Windows")
	}
	os.Args = []string{"catnet", "scan", "10.0.0.1-255", "--ping-timeout", "1000"}

	// Trigger a background cancel
	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		p.Signal(os.Interrupt)
	}()

	err := cli.Execute()
	
	if err == nil {
		t.Fatalf("Expected error due to cancellation")
	}
	
	if exitErr, ok := err.(*cli.ExitError); !ok || exitErr.Code != cli.ExitCodeInterrupted {
		t.Errorf("Expected ExitCodeInterrupted, got %v", err)
	}
}

func TestScanInvalidTarget(t *testing.T) {
	os.Args = []string{"catnet", "scan", "not-a-valid-ip"}
	
	err := cli.Execute()
	if err == nil {
		t.Fatalf("Expected error for invalid target")
	}
	
	if exitErr, ok := err.(*cli.ExitError); !ok || exitErr.Code != cli.ExitCodeInputError {
		t.Errorf("Expected ExitCodeInputError, got %v", err)
	}
}

func TestExportCSVFromJSON(t *testing.T) {
	// Create a temp JSON file
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "input.json")
	
	jsonBytes, err := os.ReadFile("../testdata/expected_output.json")
	if err != nil {
		t.Fatalf("Failed to read testdata: %v", err)
	}
	os.WriteFile(jsonPath, jsonBytes, 0644)

	os.Args = []string{"catnet", "export", jsonPath, "--format", "csv"}
	
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = cli.Execute()
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	
	outStr := buf.String()
	if !bytes.Contains([]byte(outStr), []byte("IP,Hostname,MAC,Status")) {
		t.Errorf("CSV output missing header, got: %s", outStr)
	}
}

func TestVersionOutput(t *testing.T) {
	os.Args = []string{"catnet", "version"}
	
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cli.Execute()
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("Expected nil error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	
	if !bytes.Contains(buf.Bytes(), []byte("catnet")) {
		t.Errorf("Expected version output to contain 'catnet', got: %s", buf.String())
	}
}
