package report

import (
	"bytes"
	"fmt"

	"github.com/orisan/review/internal/model"
)

type Markdown struct{}

func (Markdown) Render(result model.ReviewResult) ([]byte, error) {
	var out bytes.Buffer
	fmt.Fprintln(&out, "# Orisan Review Report")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Summary")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "| Field | Value |")
	fmt.Fprintln(&out, "|---|---|")
	fmt.Fprintf(&out, "| Decision | %s |\n", result.Summary.Decision)
	fmt.Fprintf(&out, "| Grade | %s |\n", result.Summary.Grade)
	fmt.Fprintf(&out, "| Changed files | %d |\n", result.Files.ChangedFiles)
	fmt.Fprintf(&out, "| Sensitive files | %d |\n", result.Files.SensitiveFiles)
	fmt.Fprintf(&out, "| Binary files skipped | %d |\n", result.Files.BinaryFilesSkipped)
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Findings")
	fmt.Fprintln(&out)
	if len(result.Findings) == 0 {
		fmt.Fprintln(&out, "No findings.")
	}
	return out.Bytes(), nil
}
