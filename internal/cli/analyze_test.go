package cli

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestAnalyzeEmptyPatchReportsNoChanges(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/unsupported/empty.patch", "--format", "text"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	for _, want := range []string{"No changes detected.", "Security review required: NO"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
}

func TestAnalyzeRejectsNonPatchText(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/unsupported/not_a_patch.txt"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitInputError {
		t.Fatalf("code = %d, want %d, stdout = %q, stderr = %q", code, ExitInputError, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "patch could not be parsed") {
		t.Fatalf("stderr missing parse error: %q", stderr.String())
	}
	if strings.Contains(stdout.String(), "Security review required: NO") {
		t.Fatalf("invalid patch produced misleading clean report: %q", stdout.String())
	}
}

func TestAnalyzeRejectsMalformedUnifiedDiff(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/unsupported/malformed_unified_diff.patch"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitInputError {
		t.Fatalf("code = %d, want %d, stdout = %q, stderr = %q", code, ExitInputError, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "patch could not be parsed") {
		t.Fatalf("stderr missing parse error: %q", stderr.String())
	}
}

func TestAnalyzeBinaryDiffRoutesHumanReview(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/unsupported/binary_file_change.patch", "--format", "text"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	for _, want := range []string{"Security review required: YES", "Binary file change detected", "Human review", "payload_stored=false"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "Binary files a/assets/logo.png and b/assets/logo.png differ") {
		t.Fatalf("binary report copied raw binary diff marker: %q", stdout.String())
	}
}

func TestAnalyzeMissingPatchPathFailsClearly(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", filepath.Join(t.TempDir(), "does-not-exist.patch")}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitInputError {
		t.Fatalf("code = %d, want %d", code, ExitInputError)
	}
	if !strings.Contains(stderr.String(), "input could not be read") {
		t.Fatalf("stderr missing clear file error: %q", stderr.String())
	}
}

func TestAnalyzeNonGitRepoFailsClearly(t *testing.T) {
	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("not git"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--repo", repo, "--base", "main", "--head", "HEAD"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitInputError {
		t.Fatalf("code = %d, want %d", code, ExitInputError)
	}
	if !strings.Contains(stderr.String(), "not a git repository") {
		t.Fatalf("stderr missing concise git error: %q", stderr.String())
	}
}
