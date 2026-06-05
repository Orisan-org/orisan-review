package rules

import "github.com/orisan/review/internal/model"

type Engine struct{}

func (Engine) Run() []model.Finding {
	return nil
}
