package report

import "github.com/orisan/review/internal/model"

type Reporter interface {
	Render(model.ReviewResult) ([]byte, error)
}
