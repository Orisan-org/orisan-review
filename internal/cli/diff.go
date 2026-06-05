package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

	gitdiff "github.com/orisan/review/internal/git"
	"github.com/orisan/review/internal/patch"
)

func runDiff(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("diff", flag.ContinueOnError)
	flags.SetOutput(stderr)

	var baseRef, headRef, repoPath string
	var staged, worktree bool
	common := defaultCommonFlags()

	flags.StringVar(&baseRef, "base", "", "base git ref")
	flags.StringVar(&headRef, "head", "", "head git ref")
	flags.BoolVar(&staged, "staged", false, "analyze staged changes")
	flags.BoolVar(&worktree, "worktree", false, "analyze working tree changes")
	flags.StringVar(&repoPath, "repo", "", "repository path")
	bindCommonFlags(flags, &common)

	if err := flags.Parse(args); err != nil {
		return ExitUsageError
	}
	if staged && worktree {
		_, _ = fmt.Fprintln(stderr, "--staged and --worktree are mutually exclusive")
		return ExitUsageError
	}
	if (baseRef == "") != (headRef == "") {
		_, _ = fmt.Fprintln(stderr, "--base and --head must be provided together")
		return ExitUsageError
	}

	mode := gitdiff.DiffModeWorktree
	if staged {
		mode = gitdiff.DiffModeStaged
	}
	if worktree {
		mode = gitdiff.DiffModeWorktree
	}
	if baseRef != "" && headRef != "" {
		mode = gitdiff.DiffModeRefs
	}

	data, err := gitdiff.CollectDiff(context.Background(), gitdiff.DiffOptions{
		RepoPath: repoPath,
		Mode:     mode,
		BaseRef:  baseRef,
		HeadRef:  headRef,
	})
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "collect diff: %v\n", err)
		return ExitInputError
	}
	doc, err := patch.Parse(data)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "parse diff: %v\n", err)
		return ExitInputError
	}
	result := resultFromPatch("git_diff", doc)
	result.Input.Source = "git_diff"
	result.Input.RepoRoot = repoPath
	result.Input.BaseRef = baseRef
	result.Input.HeadRef = headRef
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
