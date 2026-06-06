package report

import (
	"bytes"
	"fmt"

	"github.com/orisan/review/internal/app"
	"github.com/orisan/review/internal/model"
)

type Terminal struct{}

func (Terminal) Render(result model.ReviewResult) ([]byte, error) {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%s %s\n", app.Name, app.Version)
	fmt.Fprintf(&out, "Input: %s\n", result.Input.Source)
	fmt.Fprintf(&out, "Changed files: %d\n", result.Files.ChangedFiles)
	fmt.Fprintf(&out, "Sensitive files: %d\n", result.Files.SensitiveFiles)
	fmt.Fprintf(&out, "Generated files: %d\n", result.Files.GeneratedFiles)
	fmt.Fprintf(&out, "Binary files skipped: %d\n", result.Files.BinaryFilesSkipped)
	fmt.Fprintf(&out, "Security review required: %s\n", yesNo(result.Summary.Decision != model.DecisionPass))
	fmt.Fprintf(&out, "Risk level: %s\n", result.Summary.RiskLevel)
	fmt.Fprintf(&out, "Review decision: %s\n", result.Summary.Decision)
	fmt.Fprintf(&out, "Grade: %s\n", result.Summary.Grade)
	if len(result.Findings) == 0 && result.Files.ChangedFiles == 0 {
		fmt.Fprintln(&out, "No changes detected.")
	}
	if len(result.Findings) > 0 {
		fmt.Fprintln(&out, "\nWhy:")
		for _, finding := range result.Findings {
			fmt.Fprintf(&out, "- %s in %s\n", finding.Title, finding.Location.Path)
		}
	}
	if len(result.Summary.Routes) > 0 {
		fmt.Fprintln(&out, "\nWho should review:")
		for _, route := range result.Summary.Routes {
			fmt.Fprintf(&out, "- %s\n", routeLabel(route))
		}
	}
	if len(result.Findings) == 0 {
		fmt.Fprintln(&out, "\nFindings: none")
		return out.Bytes(), nil
	}
	fmt.Fprintln(&out, "\nFindings:")
	for i, finding := range result.Findings {
		fmt.Fprintf(&out, "%d. %s\n", i+1, finding.Category)
		fmt.Fprintf(&out, "   File: %s\n", finding.Location.Path)
		if finding.Location.StartLine != 0 {
			fmt.Fprintf(&out, "   Line hint: %d\n", finding.Location.StartLine)
		}
		fmt.Fprintf(&out, "   Reason: %s\n", finding.Evidence)
		fmt.Fprintf(&out, "   Evidence: redacted snippet only\n")
		fmt.Fprintf(&out, "   payload_stored=%v\n", finding.PayloadStored)
	}
	return out.Bytes(), nil
}

func yesNo(v bool) string {
	if v {
		return "YES"
	}
	return "NO"
}

func routeLabel(route model.ReviewRoute) string {
	switch route {
	case model.RouteAppSec:
		return "AppSec"
	case model.RouteCICD:
		return "CI/CD owner"
	case model.RouteDependency:
		return "Dependency owner"
	case model.RouteInfra:
		return "Infra owner"
	case model.RouteData:
		return "Data owner"
	case model.RouteHuman:
		return "Human review"
	default:
		return string(route)
	}
}
