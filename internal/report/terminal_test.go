package report

import (
	"strings"
	"testing"

	"github.com/orisan/review/internal/model"
)

func TestTerminalFindingReasonAndSafeEvidence(t *testing.T) {
	result := terminalResult([]model.Finding{{
		ID:            "REVIEW-AUTH-011",
		Title:         "Authorization check removed or weakened",
		Severity:      model.SeverityCritical,
		Category:      "authorization_weakened",
		Location:      model.Location{Path: "internal/api/users.go", StartLine: 42},
		Evidence:      "return nil",
		PayloadStored: false,
	}})

	out, err := (Terminal{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, want := range []string{
		"Security review required: YES",
		"Risk level: HIGH",
		"Reason: Authorization logic appears weakened by the changed condition",
		"Safe evidence: return nil",
		"payload_stored=false",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("terminal output missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "Reason: return nil") {
		t.Fatalf("terminal reason used snippet instead of explanation:\n%s", text)
	}
}

func TestTerminalDoesNotLeakSecretEvidence(t *testing.T) {
	result := terminalResult([]model.Finding{{
		ID:            "REVIEW-SEC-040",
		Title:         "Secret-like value added",
		Severity:      model.SeverityCritical,
		Category:      "secret_like_value_added",
		Location:      model.Location{Path: "config/prod.env", StartLine: 2},
		Evidence:      `AWS_SECRET_ACCESS_KEY="REDACTED"`,
		PayloadStored: false,
	}})

	out, err := (Terminal{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	if strings.Contains(text, "wJalrXUtnFEMI") || strings.Contains(text, "supersecretpassword") {
		t.Fatalf("terminal output leaked secret:\n%s", text)
	}
	if !strings.Contains(text, "Safe evidence: AWS_SECRET_ACCESS_KEY=\"REDACTED\"") {
		t.Fatalf("terminal output missing redacted evidence:\n%s", text)
	}
}

func TestTerminalSafeReadmeRemainsConcise(t *testing.T) {
	result := terminalResult(nil)
	result.Summary.Decision = model.DecisionPass
	result.Summary.RiskLevel = "NONE"
	result.Summary.Grade = "A"
	result.Files.ChangedFiles = 1
	result.Files.SensitiveFiles = 1

	out, err := (Terminal{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	if !strings.Contains(text, "Security review required: NO") || !strings.Contains(text, "Findings: none") {
		t.Fatalf("safe terminal output not concise:\n%s", text)
	}
}

func TestTerminalBinaryDiffRoutesHumanReview(t *testing.T) {
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

	out, err := (Terminal{}).Render(result)
	if err != nil {
		t.Fatal(err)
	}
	text := string(out)
	for _, want := range []string{"Human review", "Binary file change cannot be inspected safely", "payload_stored=false"} {
		if !strings.Contains(text, want) {
			t.Fatalf("binary terminal output missing %q:\n%s", want, text)
		}
	}
}

func terminalResult(findings []model.Finding) model.ReviewResult {
	decision := model.DecisionPass
	risk := "NONE"
	grade := "A"
	routes := []model.ReviewRoute{}
	if len(findings) > 0 {
		decision = model.DecisionSecurityReview
		risk = "HIGH"
		grade = "D"
		routes = []model.ReviewRoute{model.RouteAppSec}
	}
	return model.ReviewResult{
		Scanner: model.ScannerInfo{Name: "orisan-review", Version: "0.1.0"},
		Input:   model.DiffInput{Source: "patch_file"},
		Summary: model.ReviewSummary{
			Decision:  decision,
			RiskLevel: risk,
			Grade:     grade,
			Routes:    routes,
		},
		Files:    model.FileSummary{ChangedFiles: len(findings)},
		Findings: findings,
	}
}
