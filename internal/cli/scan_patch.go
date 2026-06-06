package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/orisan/review/internal/app"
	"github.com/orisan/review/internal/classify"
	"github.com/orisan/review/internal/model"
	"github.com/orisan/review/internal/patch"
	"github.com/orisan/review/internal/report"
	"github.com/orisan/review/internal/route"
	"github.com/orisan/review/internal/rules"
	"github.com/orisan/review/internal/scoring"
)

func runScanPatch(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("scan-patch", flag.ContinueOnError)
	flags.SetOutput(stderr)

	common := defaultCommonFlags()
	bindCommonFlags(flags, &common)

	if err := flags.Parse(args); err != nil {
		return ExitUsageError
	}
	if flags.NArg() != 1 {
		_, _ = fmt.Fprintln(stderr, "scan-patch requires exactly one patch path or '-' for stdin")
		return ExitUsageError
	}

	path := flags.Arg(0)
	data, err := readPatchInput(path, stdin, common.maxPatchBytes)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "read patch: %v\n", err)
		return ExitInputError
	}

	doc, err := patch.Parse(data)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "parse patch: %v\n", err)
		return ExitInputError
	}

	result := resultFromPatch(path, doc)
	rendered, err := renderResult(common.output, result)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "render report: %v\n", err)
		return ExitUsageError
	}

	if common.outPath != "" {
		if err := writeReport(common.outPath, rendered); err != nil {
			_, _ = fmt.Fprintf(stderr, "write report: %v\n", err)
			return ExitInputError
		}
		return ExitOK
	}

	_, _ = stdout.Write(rendered)
	if len(rendered) == 0 || rendered[len(rendered)-1] != '\n' {
		_, _ = fmt.Fprintln(stdout)
	}
	return ExitOK
}

func readPatchInput(path string, stdin io.Reader, maxBytes int64) ([]byte, error) {
	limited := maxBytes + 1
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(io.LimitReader(stdin, limited))
	} else {
		file, openErr := os.Open(filepath.Clean(path))
		if openErr != nil {
			return nil, openErr
		}
		defer file.Close()
		data, err = io.ReadAll(io.LimitReader(file, limited))
	}
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("patch exceeds max size of %d bytes", maxBytes)
	}
	return data, nil
}

func resultFromPatch(path string, doc patch.Document) model.ReviewResult {
	files := doc.Files
	sensitiveFiles := 0
	generatedFiles := 0
	binaryFiles := 0
	for i := range files {
		files[i].Language = classify.LanguageForPath(files[i].NewPath)
		files[i].IsGenerated = classify.IsGeneratedPath(files[i].NewPath)
		if files[i].IsGenerated {
			generatedFiles++
		}
		if files[i].IsBinary {
			binaryFiles++
		}
		if classify.IsSensitivePath(files[i].NewPath) {
			sensitiveFiles++
		}
	}

	source := "patch_file"
	if path == "-" {
		source = "stdin"
	}
	findings := (rules.Engine{}).Run(files)

	return model.ReviewResult{
		Scanner: model.ScannerInfo{
			Name:    app.Name,
			Version: app.Version,
		},
		Input: model.DiffInput{
			Source: source,
		},
		Summary: model.ReviewSummary{
			Decision:  scoring.Decision(findings),
			Routes:    route.RoutesForFindings(findings),
			Grade:     scoring.Grade(findings),
			RiskLevel: scoring.RiskLevel(findings),
			Counts:    scoring.Counts(findings),
		},
		Files: model.FileSummary{
			ChangedFiles:       len(files),
			SensitiveFiles:     sensitiveFiles,
			GeneratedFiles:     generatedFiles,
			BinaryFilesSkipped: binaryFiles,
		},
		Findings: findings,
	}
}

func renderResult(format string, result model.ReviewResult) ([]byte, error) {
	switch strings.ToLower(format) {
	case "table", "terminal", "text":
		return (report.Terminal{}).Render(result)
	case "json":
		return (report.JSON{}).Render(result)
	case "md", "markdown":
		return (report.Markdown{}).Render(result)
	case "sarif":
		return (report.SARIF{}).Render(result)
	case "html":
		return (report.HTML{}).Render(result)
	default:
		return nil, fmt.Errorf("unsupported output format %q", format)
	}
}

func writeReport(path string, data []byte) error {
	return os.WriteFile(path, data, 0o600)
}
