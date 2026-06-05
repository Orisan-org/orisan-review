package cli

import (
	"fmt"
	"io"

	"github.com/orisan/review/internal/app"
)

func runVersion(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		_, _ = fmt.Fprintln(stderr, "version does not accept arguments")
		return ExitUsageError
	}
	_, _ = fmt.Fprintf(stdout, "%s %s\n", app.Name, app.Version)
	return ExitOK
}

func bindCommonFlags(flags interface {
	StringVar(*string, string, string, string)
	Int64Var(*int64, string, int64, string)
	BoolVar(*bool, string, bool, string)
}, common *commonFlags) {
	flags.StringVar(&common.output, "output", common.output, "output format: table, json, md, sarif")
	flags.StringVar(&common.outPath, "out", "", "write report to path")
	flags.StringVar(&common.severityThreshold, "severity-threshold", common.severityThreshold, "severity threshold: low, medium, high, critical")
	flags.StringVar(&common.configPath, "config", "", "optional review.yaml path")
	flags.Int64Var(&common.maxPatchBytes, "max-patch-bytes", common.maxPatchBytes, "maximum patch size in bytes")
	flags.Int64Var(&common.maxFileBytes, "max-file-bytes", common.maxFileBytes, "maximum file size in bytes")
	flags.BoolVar(&common.noColor, "no-color", false, "disable color output")
	flags.BoolVar(&common.failOnReview, "fail-on-review-required", false, "exit nonzero when review is required")
	flags.BoolVar(&common.includeGenerated, "include-generated", false, "include generated files")
}
