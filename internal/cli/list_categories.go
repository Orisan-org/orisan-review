package cli

import (
	"fmt"
	"io"

	"github.com/orisan/review/internal/classify"
)

func runListCategories(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		_, _ = fmt.Fprintln(stderr, "list-categories does not accept arguments")
		return ExitUsageError
	}
	for _, category := range classify.Categories() {
		_, _ = fmt.Fprintln(stdout, category)
	}
	return ExitOK
}
