package report

import (
	"bytes"
	"fmt"

	"github.com/orisan/review/internal/model"
)

type Markdown struct{}

func (Markdown) Render(result model.ReviewResult) ([]byte, error) {
	var out bytes.Buffer
	fmt.Fprintln(&out, "# Orisan Review Report")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Summary")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "| Field | Value |")
	fmt.Fprintln(&out, "|---|---|")
	fmt.Fprintf(&out, "| Decision | %s |\n", result.Summary.Decision)
	fmt.Fprintf(&out, "| Risk level | %s |\n", result.Summary.RiskLevel)
	fmt.Fprintf(&out, "| Grade | %s |\n", result.Summary.Grade)
	fmt.Fprintf(&out, "| Changed files | %d |\n", result.Files.ChangedFiles)
	fmt.Fprintf(&out, "| Sensitive files | %d |\n", result.Files.SensitiveFiles)
	fmt.Fprintf(&out, "| Binary files skipped | %d |\n", result.Files.BinaryFilesSkipped)
	fmt.Fprintf(&out, "| Security review required | %t |\n", result.Summary.Decision != model.DecisionPass)
	fmt.Fprintln(&out)
	if len(result.Summary.Routes) > 0 {
		fmt.Fprintln(&out, "## Required Review")
		fmt.Fprintln(&out)
		for _, route := range result.Summary.Routes {
			fmt.Fprintf(&out, "- %s\n", route)
		}
		fmt.Fprintln(&out)
	}
	fmt.Fprintln(&out, "## Findings")
	fmt.Fprintln(&out)
	if len(result.Findings) == 0 {
		fmt.Fprintln(&out, "No findings.")
		return out.Bytes(), nil
	}
	for _, finding := range result.Findings {
		fmt.Fprintf(&out, "### %s - %s\n\n", finding.ID, finding.Title)
		fmt.Fprintf(&out, "Severity: %s  \n", finding.Severity)
		fmt.Fprintf(&out, "Category: %s  \n", finding.Category)
		fmt.Fprintf(&out, "Location: `%s:%d`  \n", finding.Location.Path, finding.Location.StartLine)
		fmt.Fprintf(&out, "payload_stored=%v  \n", finding.PayloadStored)
		fmt.Fprintf(&out, "Evidence: %s\n\n", finding.Evidence)
	}
	return out.Bytes(), nil
}
