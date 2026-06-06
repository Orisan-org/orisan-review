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
			fmt.Fprintf(&out, "- %s in %s.\n", reasonSentence(finding), finding.Location.Path)
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
		fmt.Fprintf(&out, "   Severity: %s\n", finding.Severity)
		fmt.Fprintf(&out, "   File: %s\n", finding.Location.Path)
		if finding.Location.StartLine != 0 {
			fmt.Fprintf(&out, "   Line: %d\n", finding.Location.StartLine)
		}
		fmt.Fprintf(&out, "   Reason: %s\n", reasonSentence(finding))
		if finding.Evidence != "" {
			fmt.Fprintf(&out, "   Safe evidence: %s\n", finding.Evidence)
		} else {
			fmt.Fprintf(&out, "   Safe evidence: redacted snippet only\n")
		}
		fmt.Fprintf(&out, "   payload_stored=%v\n", finding.PayloadStored)
	}
	return out.Bytes(), nil
}

func reasonSentence(finding model.Finding) string {
	switch finding.Category {
	case "authorization_weakened":
		return "Authorization logic appears weakened by the changed condition"
	case "auth_logic_changed":
		return "Authentication-sensitive logic changed"
	case "validation_removed":
		return "Input validation appears removed or bypassed"
	case "tls_verification_disabled":
		return "TLS verification appears disabled"
	case "secret_like_value_added":
		return "A secret-like value appears in added code"
	case "ci_permissions_broadened":
		return "CI workflow permissions appear broadened"
	case "unpinned_github_action":
		return "A third-party GitHub Action is not pinned to a full commit SHA"
	case "dependency_manifest_changed":
		return "A dependency manifest or lockfile changed"
	case "infra_public_exposure":
		return "Infrastructure appears to expose a service or network path publicly"
	case "destructive_migration":
		return "A migration contains a destructive database operation"
	case "tests_skipped":
		return "Tests appear skipped or disabled"
	case "ai_generated_marker":
		return "An AI-assistance marker appears in changed content"
	case "binary_file_change":
		return "Binary file change cannot be inspected safely"
	default:
		if finding.Impact != "" {
			return finding.Impact
		}
		return finding.Title
	}
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
