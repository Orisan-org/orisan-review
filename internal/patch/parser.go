package patch

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/orisan/review/internal/model"
)

var hunkHeaderPattern = regexp.MustCompile(`^@@ -([0-9]+)(?:,([0-9]+))? \+([0-9]+)(?:,([0-9]+))? @@`)

func Parse(input []byte) (Document, error) {
	text := strings.ReplaceAll(string(input), "\r\n", "\n")
	lines := strings.Split(text, "\n")

	var doc Document
	var current *model.ChangedFile
	var currentHunk *model.DiffHunk
	var oldLine, newLine int

	flushHunk := func() {
		if current != nil && currentHunk != nil {
			current.Hunks = append(current.Hunks, *currentHunk)
			currentHunk = nil
		}
	}
	flushFile := func() {
		flushHunk()
		if current != nil {
			normalizeStatus(current)
			doc.Files = append(doc.Files, *current)
			current = nil
		}
	}

	for i, line := range lines {
		if line == "" && i == len(lines)-1 {
			break
		}

		if strings.HasPrefix(line, "diff --git ") {
			flushFile()
			oldPath, newPath := parseDiffGitPaths(line)
			current = &model.ChangedFile{
				OldPath: oldPath,
				NewPath: newPath,
				Status:  model.FileModified,
			}
			continue
		}
		if current == nil {
			continue
		}

		switch {
		case strings.HasPrefix(line, "new file mode "):
			current.Status = model.FileAdded
		case strings.HasPrefix(line, "deleted file mode "):
			current.Status = model.FileDeleted
		case strings.HasPrefix(line, "rename from "):
			current.Status = model.FileRenamed
			current.OldPath = strings.TrimPrefix(line, "rename from ")
		case strings.HasPrefix(line, "rename to "):
			current.Status = model.FileRenamed
			current.NewPath = strings.TrimPrefix(line, "rename to ")
		case strings.HasPrefix(line, "Binary files "):
			current.IsBinary = true
		case strings.HasPrefix(line, "--- "):
			current.OldPath = cleanPatchPath(strings.TrimPrefix(line, "--- "))
		case strings.HasPrefix(line, "+++ "):
			current.NewPath = cleanPatchPath(strings.TrimPrefix(line, "+++ "))
		case strings.HasPrefix(line, "@@ "):
			flushHunk()
			hunk, parsedOldLine, parsedNewLine, err := parseHunkHeader(line)
			if err != nil {
				return Document{}, fmt.Errorf("parse hunk header: %w", err)
			}
			currentHunk = &hunk
			oldLine = parsedOldLine
			newLine = parsedNewLine
		case currentHunk != nil:
			switch {
			case strings.HasPrefix(line, `\ No newline at end of file`):
				continue
			case strings.HasPrefix(line, "+"):
				currentHunk.Lines = append(currentHunk.Lines, model.DiffLine{
					Type:    model.LineAdded,
					NewLine: newLine,
					Content: line[1:],
				})
				newLine++
			case strings.HasPrefix(line, "-"):
				currentHunk.Lines = append(currentHunk.Lines, model.DiffLine{
					Type:    model.LineRemoved,
					OldLine: oldLine,
					Content: line[1:],
				})
				oldLine++
			case strings.HasPrefix(line, " "):
				currentHunk.Lines = append(currentHunk.Lines, model.DiffLine{
					Type:    model.LineContext,
					OldLine: oldLine,
					NewLine: newLine,
					Content: line[1:],
				})
				oldLine++
				newLine++
			default:
				return Document{}, fmt.Errorf("unexpected diff line %q", line)
			}
		}
	}

	flushFile()
	return doc, nil
}

func parseDiffGitPaths(line string) (string, string) {
	fields := strings.Fields(strings.TrimPrefix(line, "diff --git "))
	if len(fields) < 2 {
		return "", ""
	}
	return cleanPatchPath(fields[0]), cleanPatchPath(fields[1])
}

func parseHunkHeader(line string) (model.DiffHunk, int, int, error) {
	matches := hunkHeaderPattern.FindStringSubmatch(line)
	if matches == nil {
		return model.DiffHunk{}, 0, 0, fmt.Errorf("invalid hunk header %q", line)
	}

	oldStart, err := strconv.Atoi(matches[1])
	if err != nil {
		return model.DiffHunk{}, 0, 0, err
	}
	newStart, err := strconv.Atoi(matches[3])
	if err != nil {
		return model.DiffHunk{}, 0, 0, err
	}

	return model.DiffHunk{
		Header:       line,
		OldStartLine: oldStart,
		NewStartLine: newStart,
	}, oldStart, newStart, nil
}

func cleanPatchPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.Trim(path, `"`)
	if path == "/dev/null" {
		return ""
	}
	if strings.HasPrefix(path, "a/") || strings.HasPrefix(path, "b/") {
		return path[2:]
	}
	return path
}

func normalizeStatus(file *model.ChangedFile) {
	if file.Status == model.FileAdded {
		file.OldPath = ""
	}
	if file.Status == model.FileDeleted {
		file.NewPath = file.OldPath
	}
	if file.NewPath == "" {
		file.NewPath = file.OldPath
	}
}
