package model

type DiffInput struct {
	RepoRoot string `json:"repo_root,omitempty"`
	BaseRef  string `json:"base_ref,omitempty"`
	HeadRef  string `json:"head_ref,omitempty"`
	Source   string `json:"source"`
}

type FileStatus string

const (
	FileAdded    FileStatus = "added"
	FileModified FileStatus = "modified"
	FileDeleted  FileStatus = "deleted"
	FileRenamed  FileStatus = "renamed"
)

type ChangedFile struct {
	OldPath     string     `json:"old_path,omitempty"`
	NewPath     string     `json:"new_path"`
	Status      FileStatus `json:"status"`
	Language    string     `json:"language,omitempty"`
	Hunks       []DiffHunk `json:"hunks,omitempty"`
	IsGenerated bool       `json:"is_generated,omitempty"`
	IsBinary    bool       `json:"is_binary,omitempty"`
}

type DiffHunk struct {
	Header       string     `json:"header"`
	OldStartLine int        `json:"old_start_line"`
	NewStartLine int        `json:"new_start_line"`
	Lines        []DiffLine `json:"lines,omitempty"`
}

type DiffLineType string

const (
	LineContext DiffLineType = "context"
	LineAdded   DiffLineType = "added"
	LineRemoved DiffLineType = "removed"
)

type DiffLine struct {
	Type    DiffLineType `json:"type"`
	OldLine int          `json:"old_line,omitempty"`
	NewLine int          `json:"new_line,omitempty"`
	Content string       `json:"content"`
}
