package model

type ScannerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type FindingCounts struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Info     int `json:"info"`
}

type ReviewSummary struct {
	Decision  ReviewDecision `json:"decision"`
	Routes    []ReviewRoute  `json:"routes"`
	Grade     string         `json:"grade"`
	RiskLevel string         `json:"risk_level"`
	Counts    FindingCounts  `json:"counts"`
}

type FileSummary struct {
	ChangedFiles       int `json:"changed_files"`
	SensitiveFiles     int `json:"sensitive_files"`
	GeneratedFiles     int `json:"generated_files"`
	BinaryFilesSkipped int `json:"binary_files_skipped"`
}

type ReviewResult struct {
	Scanner  ScannerInfo   `json:"scanner"`
	Input    DiffInput     `json:"input"`
	Summary  ReviewSummary `json:"summary"`
	Files    FileSummary   `json:"files"`
	Findings []Finding     `json:"findings"`
	Warnings []string      `json:"warnings,omitempty"`
}
