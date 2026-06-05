package report

import (
	"bytes"
	"fmt"

	"github.com/orisan/review/internal/app"
	"github.com/orisan/review/internal/model"
)

type Terminal struct{}

func (Terminal) Render(result model.ReviewResult) ([]byte, error) {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%s %s\n", app.Name, app.Version)
	fmt.Fprintf(&out, "Input: %s\n", result.Input.Source)
	fmt.Fprintf(&out, "Changed files: %d\n", result.Files.ChangedFiles)
	fmt.Fprintf(&out, "Sensitive files: %d\n", result.Files.SensitiveFiles)
	fmt.Fprintf(&out, "Generated files: %d\n", result.Files.GeneratedFiles)
	fmt.Fprintf(&out, "Binary files skipped: %d\n", result.Files.BinaryFilesSkipped)
	fmt.Fprintf(&out, "Review decision: %s\n", result.Summary.Decision)
	fmt.Fprintf(&out, "Grade: %s\n", result.Summary.Grade)
	if len(result.Summary.Routes) > 0 {
		fmt.Fprintln(&out, "Required reviewers:")
		for _, route := range result.Summary.Routes {
			fmt.Fprintf(&out, "- %s\n", route)
		}
	}
	if len(result.Findings) == 0 {
		fmt.Fprintln(&out, "Findings: none")
	}
	return out.Bytes(), nil
}
