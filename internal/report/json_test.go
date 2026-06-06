package report

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/orisan/review/internal/model"
)

func TestJSONRender(t *testing.T) {
	out, err := (JSON{}).Render(model.ReviewResult{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !json.Valid(out) {
		t.Fatalf("invalid json: %s", out)
	}
}

func TestSARIFDoesNotIncludeEvidence(t *testing.T) {
	result := model.ReviewResult{
		Scanner: model.ScannerInfo{Name: "orisan-review", Version: "0.1.0"},
		Findings: []model.Finding{{
			ID:            "REVIEW-SEC-040",
			Title:         "Secret-like value added",
			Severity:      model.SeverityCritical,
			Category:      "secret_like_value_added",
			Location:      model.Location{Path: "config/prod.env", StartLine: 2},
			Evidence:      `AWS_SECRET_ACCESS_KEY="wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"`,
			Remediation:   "Remove the secret.",
			PayloadStored: false,
		}},
	}
	out, err := (SARIF{}).Render(result)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !json.Valid(out) {
		t.Fatalf("invalid sarif json: %s", out)
	}
	text := string(out)
	if strings.Contains(text, "wJalrXUtnFEMI") || strings.Contains(text, "AWS_SECRET_ACCESS_KEY") {
		t.Fatalf("SARIF leaked evidence: %s", text)
	}
	if !strings.Contains(text, "payload_stored=false") {
		t.Fatalf("SARIF missing safety signal: %s", text)
	}
}
