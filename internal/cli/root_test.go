package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"--help"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d", code, ExitOK)
	}
	if !strings.Contains(stdout.String(), "orisan-review") {
		t.Fatalf("help output missing binary name: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"unknown"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitUsageError {
		t.Fatalf("code = %d, want %d", code, ExitUsageError)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("stderr missing unknown command message: %q", stderr.String())
	}
}
