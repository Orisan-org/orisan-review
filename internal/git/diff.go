package git

import (
	"context"
	"fmt"
	"os/exec"
)

type DiffMode string

const (
	DiffModeRefs     DiffMode = "refs"
	DiffModeStaged   DiffMode = "staged"
	DiffModeWorktree DiffMode = "worktree"
)

type DiffOptions struct {
	RepoPath string
	Mode     DiffMode
	BaseRef  string
	HeadRef  string
}

func BuildDiffArgs(options DiffOptions) ([]string, error) {
	args := []string{"diff", "--unified=80", "--no-ext-diff", "--no-color"}
	switch options.Mode {
	case DiffModeRefs:
		if options.BaseRef == "" || options.HeadRef == "" {
			return nil, fmt.Errorf("base and head refs are required")
		}
		args = append(args, options.BaseRef+"..."+options.HeadRef)
	case DiffModeStaged:
		args = append(args, "--cached")
	case DiffModeWorktree, "":
	default:
		return nil, fmt.Errorf("unsupported diff mode %q", options.Mode)
	}
	return args, nil
}

func CollectDiff(ctx context.Context, options DiffOptions) ([]byte, error) {
	args, err := BuildDiffArgs(options)
	if err != nil {
		return nil, err
	}
	if err := ensureGitRepo(ctx, options.RepoPath); err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "git", args...)
	if options.RepoPath != "" {
		cmd.Dir = options.RepoPath
	}
	out, err := cmd.Output()
	if err == nil {
		return out, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return nil, fmt.Errorf("git diff failed: %s", string(exitErr.Stderr))
	}
	return nil, err
}

func ensureGitRepo(ctx context.Context, repoPath string) error {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--is-inside-work-tree")
	if repoPath != "" {
		cmd.Dir = repoPath
	}
	out, err := cmd.Output()
	if err != nil {
		if repoPath == "" {
			return fmt.Errorf("not a git repository")
		}
		return fmt.Errorf("not a git repository: %s", repoPath)
	}
	if string(out) != "true\n" {
		if repoPath == "" {
			return fmt.Errorf("not a git repository")
		}
		return fmt.Errorf("not a git repository: %s", repoPath)
	}
	return nil
}
