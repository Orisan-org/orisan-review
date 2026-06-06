package git

import (
	"context"
	"testing"
)

func TestRefRange(t *testing.T) {
	r := RefRange{Base: "main", Head: "HEAD"}
	if r.Base == "" || r.Head == "" {
		t.Fatal("expected refs to be stored")
	}
}

func TestEnsureGitRepoRejectsNonRepo(t *testing.T) {
	err := ensureGitRepo(context.Background(), t.TempDir())
	if err == nil {
		t.Fatal("expected non-git repo error")
	}
}

func TestBuildDiffArgs(t *testing.T) {
	cases := []struct {
		name    string
		options DiffOptions
		want    []string
	}{
		{
			name:    "worktree",
			options: DiffOptions{Mode: DiffModeWorktree},
			want:    []string{"diff", "--unified=80", "--no-ext-diff", "--no-color"},
		},
		{
			name:    "staged",
			options: DiffOptions{Mode: DiffModeStaged},
			want:    []string{"diff", "--unified=80", "--no-ext-diff", "--no-color", "--cached"},
		},
		{
			name:    "refs",
			options: DiffOptions{Mode: DiffModeRefs, BaseRef: "main", HeadRef: "HEAD"},
			want:    []string{"diff", "--unified=80", "--no-ext-diff", "--no-color", "main...HEAD"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := BuildDiffArgs(tc.options)
			if err != nil {
				t.Fatalf("BuildDiffArgs() error = %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("args = %v, want %v", got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("args = %v, want %v", got, tc.want)
				}
			}
		})
	}
}
