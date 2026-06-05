package redact

import "testing"

func TestEvidenceScaffold(t *testing.T) {
	if Evidence("hello") != "hello" {
		t.Fatal("expected scaffold redaction to preserve ordinary text")
	}
}
