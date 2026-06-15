package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunScanPatchFromStdin(t *testing.T) {
	input := `diff --git a/README.md b/README.md
index 1111111..2222222 100644
--- a/README.md
+++ b/README.md
@@ -1 +1 @@
-Hello
+Hello.
`
	var stdout, stderr bytes.Buffer
	code := Run([]string{"scan-patch", "--output", "json", "-"}, strings.NewReader(input), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"changed_files": 1`) {
		t.Fatalf("stdout missing changed file count: %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"decision": "pass"`) {
		t.Fatalf("stdout missing pass decision: %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"sensitive_files": 0`) {
		t.Fatalf("README-only diff should not count as sensitive: %s", stdout.String())
	}
}
