package redact

import (
	"strings"
	"testing"
)

func TestEvidenceScaffold(t *testing.T) {
	if Evidence("hello") != "hello" {
		t.Fatal("expected scaffold redaction to preserve ordinary text")
	}
}

func TestEvidenceRedactsSecrets(t *testing.T) {
	got := Evidence(`AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"`)
	if got == `AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"` {
		t.Fatal("expected secret to be redacted")
	}
	if got != `AWS_SECRET_ACCESS_KEY="REDACTED"` {
		t.Fatalf("Evidence() = %q", got)
	}
}

func TestEvidenceRedactsGitHubToken(t *testing.T) {
	got := Evidence(`Authorization: Bearer ghp_abcdefghijklmnopqrstuvwxyz123456`)
	if strings.Contains(got, "abcdefghijklmnopqrstuvwxyz") {
		t.Fatalf("Evidence leaked token: %q", got)
	}
}
