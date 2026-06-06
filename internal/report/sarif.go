package report

import (
	"encoding/json"

	"github.com/orisan/review/internal/model"
)

type SARIF struct{}

func (SARIF) Render(result model.ReviewResult) ([]byte, error) {
	report := sarifLog{
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{Driver: sarifDriver{
				Name:            result.Scanner.Name,
				SemanticVersion: result.Scanner.Version,
				Rules:           sarifRules(result.Findings),
			}},
			Results: sarifResults(result.Findings),
		}},
	}
	return json.MarshalIndent(report, "", "  ")
}

type sarifLog struct {
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name            string      `json:"name"`
	SemanticVersion string      `json:"semanticVersion"`
	Rules           []sarifRule `json:"rules,omitempty"`
}

type sarifRule struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	ShortDesc sarifMessage `json:"shortDescription"`
	Help      sarifMessage `json:"help,omitempty"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations,omitempty"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

func sarifRules(findings []model.Finding) []sarifRule {
	seen := map[string]bool{}
	var rules []sarifRule
	for _, finding := range findings {
		if seen[finding.ID] {
			continue
		}
		seen[finding.ID] = true
		rules = append(rules, sarifRule{
			ID:        finding.ID,
			Name:      finding.Title,
			ShortDesc: sarifMessage{Text: finding.Title},
			Help:      sarifMessage{Text: finding.Remediation},
		})
	}
	return rules
}

func sarifResults(findings []model.Finding) []sarifResult {
	results := make([]sarifResult, 0, len(findings))
	for _, finding := range findings {
		result := sarifResult{
			RuleID:  finding.ID,
			Level:   sarifLevel(finding.Severity),
			Message: sarifMessage{Text: finding.Title + "; payload_stored=false"},
		}
		if finding.Location.Path != "" {
			result.Locations = []sarifLocation{{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: finding.Location.Path},
					Region:           sarifRegion{StartLine: finding.Location.StartLine},
				},
			}}
		}
		results = append(results, result)
	}
	return results
}

func sarifLevel(severity model.Severity) string {
	switch severity {
	case model.SeverityCritical, model.SeverityHigh:
		return "error"
	case model.SeverityMedium:
		return "warning"
	default:
		return "note"
	}
}
