package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mendsec/catnet/internal/cli"
)

var binaryPath string

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	// Build the binary
	tmpDir, err := os.MkdirTemp("", "catnet-test")
	if err != nil {
		return 1
	}
	defer os.RemoveAll(tmpDir)

	if runtime.GOOS == "windows" {
		binaryPath = filepath.Join(tmpDir, "catnet.exe")
	} else {
		binaryPath = filepath.Join(tmpDir, "catnet")
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, "../cmd/catnet")
	if err := cmd.Run(); err != nil {
		return 1
	}

	return m.Run()
}

func startTestServer(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	t.Cleanup(func() { ln.Close() })
	return port
}

func TestScanOutputJSON(t *testing.T) {
	port := startTestServer(t)
	portStr := fmt.Sprintf("%d", port)

	cmd := exec.Command(binaryPath, "scan", "127.0.0.1", "--format", "json", "--quiet", "--ports", portStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected nil error, got %v: %s", err, out)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(out, &data); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput was: %s", err, out)
	}
	
	if ver, ok := data["schemaVersion"].(string); !ok || ver != "2.0.0" {
		t.Errorf("Expected schemaVersion 2.0.0, got %v", ver)
	}
	if alive, ok := data["alive"].(float64); !ok || alive != 1 {
		t.Errorf("Expected 1 alive host, got %v", alive)
	}
}

func TestScanCancelledByContext(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Signal testing is unreliable on Windows")
	}
	cmd := exec.Command(binaryPath, "scan", "10.0.0.1-10.0.255.255", "--ping-timeout", "1000")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}

	// Trigger a background cancel
	go func() {
		time.Sleep(1 * time.Second)
		cmd.Process.Signal(os.Interrupt)
	}()

	err := cmd.Wait()
	
	if err == nil {
		t.Fatalf("Expected error due to cancellation")
	}
	
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != cli.ExitCodeInterrupted {
			t.Errorf("Expected ExitCodeInterrupted (%d), got %v\nStderr: %s", cli.ExitCodeInterrupted, exitErr.ExitCode(), stderr.String())
		}
	} else {
		t.Errorf("Expected ExitError, got %v\nStderr: %s", err, stderr.String())
	}
}

func TestScanInvalidTarget(t *testing.T) {
	cmd := exec.Command(binaryPath, "scan", "not-a-valid-ip")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error for invalid target")
	}
	
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != cli.ExitCodeInputError {
			t.Errorf("Expected ExitCodeInputError (%d), got %v", cli.ExitCodeInputError, exitErr.ExitCode())
		}
	} else {
		t.Errorf("Expected ExitError, got %v", err)
	}
}

func TestExportXMLFromJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "input.json")

	jsonBytes, err := os.ReadFile("../testdata/expected_output.json")
	if err != nil {
		t.Fatalf("Failed to read testdata: %v", err)
	}
	os.WriteFile(jsonPath, jsonBytes, 0644)

	cmd := exec.Command(binaryPath, "export", jsonPath, "--format", "xml")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected nil error, got %v: %s", err, out)
	}

	if !bytes.Contains(out, []byte("<results>")) {
		t.Errorf("XML output missing root element, got: %s", out)
	}
	if !bytes.Contains(out, []byte("<ip>127.0.0.1</ip>")) {
		t.Errorf("XML output missing device IP, got: %s", out)
	}
	if !bytes.Contains(out, []byte("<status>Alive</status>")) {
		t.Errorf("XML output missing status element, got: %s", out)
	}
}

func TestExportCSVFromJSON(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "input.json")
	
	jsonBytes, err := os.ReadFile("../testdata/expected_output.json")
	if err != nil {
		t.Fatalf("Failed to read testdata: %v", err)
	}
	os.WriteFile(jsonPath, jsonBytes, 0644)

	cmd := exec.Command(binaryPath, "export", jsonPath, "--format", "csv")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected nil error, got %v: %s", err, out)
	}

	if !bytes.Contains(out, []byte("IP,Hostname,MAC,Status")) {
		t.Errorf("CSV output missing header, got: %s", out)
	}
}

func TestExportInvalidInputFile(t *testing.T) {
	cmd := exec.Command(binaryPath, "export", "/nonexistent/path/input.json", "--format", "json")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error for non-existent input file")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != cli.ExitCodeInputError {
			t.Errorf("Expected ExitCodeInputError (%d), got %v", cli.ExitCodeInputError, exitErr.ExitCode())
		}
	}
}

func TestExportInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	badPath := filepath.Join(tmpDir, "bad.json")
	os.WriteFile(badPath, []byte("{invalid json}"), 0644)

	cmd := exec.Command(binaryPath, "export", badPath, "--format", "json")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error for invalid JSON")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != cli.ExitCodeInputError {
			t.Errorf("Expected ExitCodeInputError (%d), got %v", cli.ExitCodeInputError, exitErr.ExitCode())
		}
	}
}

func TestExportUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "input.json")
	jsonBytes, _ := os.ReadFile("../testdata/expected_output.json")
	os.WriteFile(jsonPath, jsonBytes, 0644)

	cmd := exec.Command(binaryPath, "export", jsonPath, "--format", "invalid")
	err := cmd.Run()
	if err == nil {
		t.Fatalf("Expected error for unsupported format")
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != cli.ExitCodeInputError {
			t.Errorf("Expected ExitCodeInputError (%d), got %v", cli.ExitCodeInputError, exitErr.ExitCode())
		}
	}
}

func TestExportSchemaVersionWarning(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "input.json")
	jsonBytes, _ := os.ReadFile("../testdata/expected_output.json")
	badVersion := string(jsonBytes)
	badVersion = strings.ReplaceAll(badVersion, "\"2.0.0\"", "\"99.0.0\"")
	os.WriteFile(jsonPath, []byte(badVersion), 0644)

	cmd := exec.Command(binaryPath, "export", jsonPath, "--format", "csv")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected nil error, got %v: %s", err, out)
	}
	if !bytes.Contains(out, []byte("IP,Hostname,MAC,Status")) {
		t.Errorf("CSV output missing header, got: %s", out)
	}
}

func TestVersionOutput(t *testing.T) {
	cmd := exec.Command(binaryPath, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected nil error, got %v: %s", err, out)
	}

	if !bytes.Contains(out, []byte("catnet")) {
		t.Errorf("Expected version output to contain 'catnet', got: %s", out)
	}
}
