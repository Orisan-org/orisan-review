package report

import (
	"encoding/json"

	"github.com/orisan/review/internal/model"
)

type JSON struct{}

func (JSON) Render(result model.ReviewResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}
