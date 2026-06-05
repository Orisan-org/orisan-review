package report

import (
	"encoding/json"
	"testing"

	"github.com/orisan/review/internal/model"
)

func TestJSONRender(t *testing.T) {
	out, err := (JSON{}).Render(model.ReviewResult{})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !json.Valid(out) {
		t.Fatalf("invalid json: %s", out)
	}
}
