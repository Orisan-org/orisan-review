package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestAnalyzeRequiresInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitUsageError {
		t.Fatalf("code = %d, want %d", code, ExitUsageError)
	}
	if !strings.Contains(stderr.String(), "error: no input provided") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestAnalyzeFromStdinText(t *testing.T) {
	input := `diff --git a/internal/http/client.go b/internal/http/client.go
index 1111111..2222222 100644
--- a/internal/http/client.go
+++ b/internal/http/client.go
@@ -1 +1 @@
-return &http.Client{}
+return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
`
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--stdin", "--format", "text"}, strings.NewReader(input), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	for _, want := range []string{"Security review required: YES", "Risk level: HIGH", "payload_stored=false"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
}
