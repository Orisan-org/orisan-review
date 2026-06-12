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

func TestAnalyzeHelpExitsOK(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--help"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Usage of analyze") {
		t.Fatalf("help output missing usage: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("help wrote to stderr: %q", stderr.String())
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
	for _, want := range []string{"Security review required: YES", "Risk level: MEDIUM", "Review decision: review_required", "Binary file change cannot be inspected safely", "Human review", "payload_stored=false"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "Binary files a/assets/logo.png and b/assets/logo.png differ") {
		t.Fatalf("binary report copied raw binary diff marker: %q", stdout.String())
	}
}

func TestAnalyzeGitBinaryPatchRoutesHumanReview(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/unsupported/git_binary_patch.patch", "--format", "text"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	for _, want := range []string{"Security review required: YES", "Risk level: MEDIUM", "Review decision: review_required", "Binary file change cannot be inspected safely", "Human review", "payload_stored=false"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "GIT binary patch") || strings.Contains(stdout.String(), "literal 4") {
		t.Fatalf("binary report copied raw git binary patch payload: %q", stdout.String())
	}
}

func TestAnalyzeGitBinaryDeltaPatchRoutesHumanReview(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/unsupported/git_binary_patch_delta.patch", "--format", "text"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	for _, want := range []string{"Security review required: YES", "Risk level: MEDIUM", "Review decision: review_required", "Binary file change cannot be inspected safely", "Human review", "payload_stored=false"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "GIT binary patch") || strings.Contains(stdout.String(), "delta 12") {
		t.Fatalf("binary report copied raw git binary delta payload: %q", stdout.String())
	}
}

func TestAnalyzeHTMLFormatAccepted(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/diffs/authorization_weakened.patch", "--format", "html"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	text := stdout.String()
	for _, want := range []string{"<!doctype html>", "Security review required", "YES", "AppSec", "payload_stored", "false"} {
		if !strings.Contains(text, want) {
			t.Fatalf("HTML stdout missing %q:\n%s", want, text)
		}
	}
}

func TestAnalyzeHTMLSafeReadmeNoFindings(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/diffs/safe_readme_change.patch", "--format", "html"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	text := stdout.String()
	for _, want := range []string{"Security review required", "NO", "No findings."} {
		if !strings.Contains(text, want) {
			t.Fatalf("safe HTML stdout missing %q:\n%s", want, text)
		}
	}
}

func TestAnalyzeHTMLBinaryDiffRoutesHumanReview(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"analyze", "--patch", "../../testdata/unsupported/binary_file_change.patch", "--format", "html"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("code = %d, want %d, stderr = %q", code, ExitOK, stderr.String())
	}
	text := stdout.String()
	for _, want := range []string{"Human review", "Binary file change cannot be inspected safely", "payload_stored"} {
		if !strings.Contains(text, want) {
			t.Fatalf("binary HTML stdout missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Binary files a/assets/logo.png and b/assets/logo.png differ") {
		t.Fatalf("binary HTML copied raw binary diff marker:\n%s", text)
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
