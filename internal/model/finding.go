package model

type Finding struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Severity      Severity          `json:"severity"`
	Category      string            `json:"category"`
	Location      Location          `json:"location"`
	Evidence      string            `json:"evidence"`
	Impact        string            `json:"impact"`
	Remediation   string            `json:"remediation"`
	Confidence    string            `json:"confidence"`
	PayloadStored bool              `json:"payload_stored"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

func NewFinding(id, title string, severity Severity) Finding {
	return Finding{
		ID:            id,
		Title:         title,
		Severity:      severity,
		PayloadStored: false,
	}
}
