package cli

import (
	"fmt"
	"io"

	"github.com/orisan/review/internal/rules"
)

func runListRules(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		_, _ = fmt.Fprintln(stderr, "list-rules does not accept arguments")
		return ExitUsageError
	}
	for _, rule := range rules.Catalogue() {
		_, _ = fmt.Fprintf(stdout, "%s\t%s\t%s\n", rule.ID, rule.Severity, rule.Title)
	}
	return ExitOK
}
