package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunDiffRejectsMixedModes(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"diff", "--staged", "--worktree"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitUsageError {
		t.Fatalf("code = %d, want %d", code, ExitUsageError)
	}
}

func TestRunDiffRequiresBaseAndHeadTogether(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"diff", "--base", "main"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitUsageError {
		t.Fatalf("code = %d, want %d", code, ExitUsageError)
	}
}
