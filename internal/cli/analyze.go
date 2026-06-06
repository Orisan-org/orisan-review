package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	gitdiff "github.com/orisan/review/internal/git"
	"github.com/orisan/review/internal/patch"
)

func runAnalyze(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("analyze", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var patchPath, repoPath, baseRef, headRef, format, outPath string
	var useStdin bool
	common := defaultCommonFlags()

	flags.StringVar(&patchPath, "patch", "", "unified git diff patch path")
	flags.BoolVar(&useStdin, "stdin", false, "read unified git diff from stdin")
	flags.StringVar(&repoPath, "repo", "", "repository path for git diff input")
	flags.StringVar(&baseRef, "base", "", "base git ref")
	flags.StringVar(&headRef, "head", "", "head git ref")
	flags.StringVar(&format, "format", "text", "output format: text, json, md, sarif, html")
	flags.StringVar(&outPath, "out", "", "write report to path")
	flags.Int64Var(&common.maxPatchBytes, "max-patch-bytes", common.maxPatchBytes, "maximum patch size in bytes")

	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			flags.SetOutput(stdout)
			break
		}
	}
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ExitOK
		}
		return ExitUsageError
	}
	if flags.NArg() > 0 {
		_, _ = fmt.Fprintf(stderr, "error: unexpected arguments: %v\n", flags.Args())
		return ExitUsageError
	}

	inputCount := 0
	if patchPath != "" {
		inputCount++
	}
	if useStdin {
		inputCount++
	}
	if repoPath != "" || baseRef != "" || headRef != "" {
		inputCount++
	}
	if inputCount == 0 {
		_, _ = fmt.Fprintln(stderr, "error: no input provided")
		_, _ = fmt.Fprintln(stderr)
		_, _ = fmt.Fprintln(stderr, "Provide one of:")
		_, _ = fmt.Fprintln(stderr, "  --patch path/to/file.patch")
		_, _ = fmt.Fprintln(stderr, "  --stdin")
		_, _ = fmt.Fprintln(stderr, "  --repo . --base main --head HEAD")
		return ExitUsageError
	}
	if inputCount > 1 {
		_, _ = fmt.Fprintln(stderr, "error: provide only one input mode")
		return ExitUsageError
	}

	var data []byte
	var source string
	var err error
	switch {
	case patchPath != "":
		data, err = readPatchInput(patchPath, stdin, common.maxPatchBytes)
		source = "patch_file"
	case useStdin:
		data, err = readPatchInput("-", stdin, common.maxPatchBytes)
		source = "stdin"
	default:
		if baseRef == "" || headRef == "" {
			_, _ = fmt.Fprintln(stderr, "error: --repo input requires --base and --head")
			return ExitUsageError
		}
		data, err = gitdiff.CollectDiff(context.Background(), gitdiff.DiffOptions{
			RepoPath: repoPath,
			Mode:     gitdiff.DiffModeRefs,
			BaseRef:  baseRef,
			HeadRef:  headRef,
		})
		source = "git_diff"
	}
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: input could not be read: %v\n", err)
		return ExitInputError
	}

	doc, err := patch.Parse(data)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "error: patch could not be parsed")
		if patchPath != "" {
			_, _ = fmt.Fprintf(stderr, "file: %s\n", patchPath)
		}
		_, _ = fmt.Fprintln(stderr, "hint: provide a unified git diff or patch file")
		return ExitInputError
	}

	result := resultFromPatch(inputLabel(patchPath, useStdin), doc)
	result.Input.Source = source
	result.Input.RepoRoot = repoPath
	result.Input.BaseRef = baseRef
	result.Input.HeadRef = headRef

	rendered, err := renderResult(format, result)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitUsageError
	}
	if outPath != "" {
		if err := os.WriteFile(outPath, rendered, 0o600); err != nil {
			_, _ = fmt.Fprintf(stderr, "error: write report: %v\n", err)
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

func inputLabel(patchPath string, useStdin bool) string {
	if useStdin {
		return "-"
	}
	if patchPath != "" {
		return patchPath
	}
	return "git_diff"
}
