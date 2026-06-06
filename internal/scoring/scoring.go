package scoring

import "github.com/orisan/review/internal/model"

func Grade(findings []model.Finding) string {
	max := maxSeverity(findings)
	switch max {
	case model.SeverityCritical:
		return "F"
	case model.SeverityHigh:
		return "D"
	case model.SeverityMedium:
		return "C"
	case model.SeverityLow, model.SeverityInfo:
		return "B"
	default:
		return "A"
	}
}

func Decision(findings []model.Finding) model.ReviewDecision {
	max := maxSeverity(findings)
	switch max {
	case model.SeverityCritical:
		return model.DecisionBlockUntilReviewed
	case model.SeverityHigh:
		return model.DecisionSecurityReview
	case model.SeverityMedium:
		return model.DecisionReviewRequired
	default:
		return model.DecisionPass
	}
}

func Counts(findings []model.Finding) model.FindingCounts {
	var counts model.FindingCounts
	for _, finding := range findings {
		switch finding.Severity {
		case model.SeverityCritical:
			counts.Critical++
		case model.SeverityHigh:
			counts.High++
		case model.SeverityMedium:
			counts.Medium++
		case model.SeverityLow:
			counts.Low++
		case model.SeverityInfo:
			counts.Info++
		}
	}
	return counts
}

func RiskLevel(findings []model.Finding) string {
	max := maxSeverity(findings)
	switch max {
	case model.SeverityCritical, model.SeverityHigh:
		return "HIGH"
	case model.SeverityMedium:
		return "MEDIUM"
	case model.SeverityLow, model.SeverityInfo:
		return "LOW"
	default:
		return "NONE"
	}
}

func maxSeverity(findings []model.Finding) model.Severity {
	max := model.Severity("")
	maxRank := -1
	for _, finding := range findings {
		if rank := model.SeverityRank(finding.Severity); rank > maxRank {
			max = finding.Severity
			maxRank = rank
		}
	}
	return max
}
