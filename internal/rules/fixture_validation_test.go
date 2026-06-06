package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/orisan/review/internal/model"
	"github.com/orisan/review/internal/patch"
	"github.com/orisan/review/internal/route"
	"github.com/orisan/review/internal/scoring"
)

type expectedFixture struct {
	SecurityReviewRequired bool     `json:"security_review_required"`
	ExpectedReviewers      []string `json:"expected_reviewers"`
	ExpectedCategories     []string `json:"expected_categories"`
	ExpectedFiles          []string `json:"expected_files"`
	PayloadStored          bool     `json:"payload_stored"`
}

func TestFixtureValidationCorpus(t *testing.T) {
	expectedFiles, err := filepath.Glob("../../testdata/expected/*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(expectedFiles) < 10 {
		t.Fatalf("expected validation corpus, got %d fixtures", len(expectedFiles))
	}

	for _, expectedPath := range expectedFiles {
		name := strings.TrimSuffix(filepath.Base(expectedPath), ".json")
		t.Run(name, func(t *testing.T) {
			expected := readExpected(t, expectedPath)
			patchBytes, err := os.ReadFile(filepath.Join("../../testdata/diffs", name+".patch"))
			if err != nil {
				t.Fatal(err)
			}
			doc, err := patch.Parse(patchBytes)
			if err != nil {
				t.Fatal(err)
			}
			findings := (Engine{}).Run(doc.Files)
			routes := route.RoutesForFindings(findings)
			decision := scoring.Decision(findings)

			if got := decision != model.DecisionPass; got != expected.SecurityReviewRequired {
				t.Fatalf("security_review_required = %v, want %v; findings = %+v", got, expected.SecurityReviewRequired, findings)
			}
			assertSet(t, "categories", findingCategories(findings), expected.ExpectedCategories)
			assertSet(t, "reviewers", reviewRoutes(routes), expected.ExpectedReviewers)
			if len(expected.ExpectedFiles) > 0 {
				assertSet(t, "files", findingFiles(findings), expected.ExpectedFiles)
			}
			for _, finding := range findings {
				if finding.PayloadStored != expected.PayloadStored {
					t.Fatalf("payload_stored = %v, want %v", finding.PayloadStored, expected.PayloadStored)
				}
				if len(finding.Evidence) > 160+3 {
					t.Fatalf("evidence too long: %d chars: %q", len(finding.Evidence), finding.Evidence)
				}
			}
		})
	}
}

func TestEvidenceLeakageCorpus(t *testing.T) {
	patchBytes, err := os.ReadFile("../../testdata/diffs/secret_like_value_added.patch")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := patch.Parse(patchBytes)
	if err != nil {
		t.Fatal(err)
	}
	findings := (Engine{}).Run(doc.Files)
	reportBytes, err := json.Marshal(findings)
	if err != nil {
		t.Fatal(err)
	}
	report := string(reportBytes)
	for _, dangerous := range []string{"wJalrXUtnFEMI", "supersecretpassword", "BEGIN PRIVATE KEY", "postgres://admin"} {
		if strings.Contains(report, dangerous) {
			t.Fatalf("report leaked dangerous substring %q: %s", dangerous, report)
		}
	}
}

func readExpected(t *testing.T, path string) expectedFixture {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var expected expectedFixture
	if err := json.Unmarshal(data, &expected); err != nil {
		t.Fatal(err)
	}
	return expected
}

func findingCategories(findings []model.Finding) []string {
	var out []string
	for _, finding := range findings {
		out = append(out, finding.Category)
	}
	return out
}

func reviewRoutes(routes []model.ReviewRoute) []string {
	var out []string
	for _, route := range routes {
		out = append(out, string(route))
	}
	return out
}

func findingFiles(findings []model.Finding) []string {
	seen := map[string]bool{}
	var out []string
	for _, finding := range findings {
		if finding.Location.Path != "" && !seen[finding.Location.Path] {
			seen[finding.Location.Path] = true
			out = append(out, finding.Location.Path)
		}
	}
	return out
}

func assertSet(t *testing.T, name string, got, want []string) {
	t.Helper()
	sort.Strings(got)
	sort.Strings(want)
	if strings.Join(got, "\x00") != strings.Join(want, "\x00") {
		t.Fatalf("%s = %v, want %v", name, got, want)
	}
}
