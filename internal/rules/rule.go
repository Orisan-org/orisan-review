package rules

import "github.com/orisan/review/internal/model"

type Metadata struct {
	ID       string
	Title    string
	Category string
	Severity model.Severity
}
