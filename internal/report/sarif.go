package report

import "github.com/orisan/review/internal/model"

type SARIF struct{}

func (SARIF) Render(_ model.ReviewResult) ([]byte, error) {
	return []byte("{}"), nil
}
