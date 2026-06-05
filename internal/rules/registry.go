package rules

import "github.com/orisan/review/internal/model"

func Catalogue() []Metadata {
	return []Metadata{
		{ID: "REVIEW-AI-001", Title: "AI-assisted change marker", Category: "ai_markers", Severity: model.SeverityInfo},
		{ID: "REVIEW-AI-002", Title: "Agent autonomy marker", Category: "ai_markers", Severity: model.SeverityMedium},
		{ID: "REVIEW-AUTH-010", Title: "Authentication logic changed", Category: "auth", Severity: model.SeverityHigh},
		{ID: "REVIEW-AUTH-011", Title: "Authorization check removed or weakened", Category: "auth", Severity: model.SeverityHigh},
		{ID: "REVIEW-VAL-020", Title: "Input validation removed or weakened", Category: "validation", Severity: model.SeverityHigh},
		{ID: "REVIEW-CRYPTO-030", Title: "TLS verification disabled", Category: "crypto", Severity: model.SeverityCritical},
		{ID: "REVIEW-CRYPTO-031", Title: "Weak crypto introduced", Category: "crypto", Severity: model.SeverityHigh},
		{ID: "REVIEW-SEC-040", Title: "Secret-like value added", Category: "secrets", Severity: model.SeverityCritical},
		{ID: "REVIEW-LOG-050", Title: "Sensitive logging introduced", Category: "logging", Severity: model.SeverityHigh},
		{ID: "REVIEW-CICD-060", Title: "GitHub Actions permissions broadened", Category: "ci_cd", Severity: model.SeverityHigh},
		{ID: "REVIEW-CICD-061", Title: "Unpinned third-party action introduced", Category: "ci_cd", Severity: model.SeverityMedium},
		{ID: "REVIEW-CICD-062", Title: "Script injection-prone workflow pattern", Category: "ci_cd", Severity: model.SeverityHigh},
		{ID: "REVIEW-DEP-070", Title: "Dependency manifest changed", Category: "dependencies", Severity: model.SeverityMedium},
		{ID: "REVIEW-DEP-071", Title: "Install script or postinstall added", Category: "dependencies", Severity: model.SeverityHigh},
		{ID: "REVIEW-INFRA-080", Title: "Public network exposure changed", Category: "infra", Severity: model.SeverityHigh},
		{ID: "REVIEW-INFRA-081", Title: "IAM/admin privilege broadened", Category: "infra", Severity: model.SeverityHigh},
		{ID: "REVIEW-DB-090", Title: "Destructive migration", Category: "data", Severity: model.SeverityHigh},
		{ID: "REVIEW-TEST-100", Title: "Tests removed or skipped", Category: "tests_bypass", Severity: model.SeverityMedium},
		{ID: "REVIEW-PUBLIC-110", Title: "Security/privacy claim changed", Category: "public_claims", Severity: model.SeverityMedium},
	}
}
