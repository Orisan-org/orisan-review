package report

import (
	"strings"
	"testing"

	"github.com/orisan/review/internal/model"
)

func TestHTMLReportFindingDetailsAndRouting(t *testing.T) {
	result := terminalResult([]model.Finding{{
		ID:            "REVIEW-AUTH-011",
		Title:         "Authorization check removed or weakened",
		Severity:      model.SeverityCritical,
		Category:      "authorization_weakened",
		Location:      model.Location{Path: "internal/api/users.go", StartLine: 42},
		Evidence:      "return nil",
		PayloadStored: false,
	}})

	out, err := (HTML{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, want := range []string{
		"Orisan Review",
		"Security review required",
		"YES",
		"Risk level",
		"Reviewer routing",
		"AppSec",
		"REVIEW-AUTH-011",
		"Authorization logic appears weakened by the changed condition",
		"Safe evidence",
		"return nil",
		"payload_stored",
		"false",
		"Generated locally",
		"No cloud calls",
		"No source upload",
		"No full diff stored",
		"Evidence is redacted",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("HTML output missing %q:\n%s", want, text)
		}
	}
}

func TestHTMLEscapesUnsafeContent(t *testing.T) {
	result := terminalResult([]model.Finding{{
		ID:            `REVIEW-XSS-001`,
		Title:         `<script>alert("title")</script>`,
		Severity:      model.SeverityHigh,
		Category:      `custom_category`,
		Location:      model.Location{Path: `internal/<api>/users.go`, StartLine: 12},
		Evidence:      `<img src=x onerror=alert("evidence")>`,
		PayloadStored: false,
	}})

	out, err := (HTML{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, raw := range []string{
		`<script>alert("title")</script>`,
		`<img src=x onerror=alert("evidence")>`,
		`internal/<api>/users.go`,
	} {
		if strings.Contains(text, raw) {
			t.Fatalf("HTML output contained unescaped raw content %q:\n%s", raw, text)
		}
	}
	for _, escaped := range []string{
		`&lt;script&gt;alert(&#34;title&#34;)&lt;/script&gt;`,
		`&lt;img src=x onerror=alert(&#34;evidence&#34;)&gt;`,
		`internal/&lt;api&gt;/users.go`,
	} {
		if !strings.Contains(text, escaped) {
			t.Fatalf("HTML output missing escaped content %q:\n%s", escaped, text)
		}
	}
}

func TestHTMLDoesNotLeakSecretEvidence(t *testing.T) {
	result := terminalResult([]model.Finding{{
		ID:            "REVIEW-SEC-040",
		Title:         "Secret-like value added",
		Severity:      model.SeverityCritical,
		Category:      "secret_like_value_added",
		Location:      model.Location{Path: "config/prod.env", StartLine: 2},
		Evidence:      `AWS_SECRET_ACCESS_KEY="REDACTED"`,
		PayloadStored: false,
	}})

	out, err := (HTML{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, leaked := range []string{"wJalrXUtnFEMI", "supersecretpassword", "BEGIN PRIVATE KEY", "postgres://admin"} {
		if strings.Contains(text, leaked) {
			t.Fatalf("HTML output leaked secret %q:\n%s", leaked, text)
		}
	}
	if !strings.Contains(text, `AWS_SECRET_ACCESS_KEY=&#34;REDACTED&#34;`) {
		t.Fatalf("HTML output missing redacted evidence:\n%s", text)
	}
}

func TestHTMLSafeReportNoFindings(t *testing.T) {
	result := terminalResult(nil)
	result.Summary.Decision = model.DecisionPass
	result.Summary.RiskLevel = "NONE"
	result.Summary.Grade = "A"
	result.Files.ChangedFiles = 1
	result.Files.SensitiveFiles = 1

	out, err := (HTML{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, want := range []string{"Security review required", "NO", "No findings.", "No reviewer route required."} {
		if !strings.Contains(text, want) {
			t.Fatalf("safe HTML output missing %q:\n%s", want, text)
		}
	}
}

func TestHTMLBinaryDiffRoutesHumanReview(t *testing.T) {
	result := terminalResult([]model.Finding{{
		ID:            "REVIEW-BIN-001",
		Title:         "Binary file change detected",
		Severity:      model.SeverityMedium,
		Category:      "binary_file_change",
		Location:      model.Location{Path: "assets/logo.png"},
		Evidence:      "Binary diff cannot be inspected safely.",
		PayloadStored: false,
	}})
	result.Summary.Routes = []model.ReviewRoute{model.RouteHuman}
	result.Files.BinaryFilesSkipped = 1

	out, err := (HTML{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, want := range []string{"Human review", "Binary file change cannot be inspected safely", "payload_stored"} {
		if !strings.Contains(text, want) {
			t.Fatalf("binary HTML output missing %q:\n%s", want, text)
		}
	}
}

func TestHTMLHasNoExternalAssetsOrScripts(t *testing.T) {
	out, err := (HTML{}).Render(terminalResult(nil))
	if err != nil {
		t.Fatal(err)
	}
	text := strings.ToLower(string(out))
	for _, forbidden := range []string{"<script", "href=\"http", "src=\"http", "@import", "url(http"} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("HTML output contains external asset or script marker %q:\n%s", forbidden, text)
		}
	}
}
