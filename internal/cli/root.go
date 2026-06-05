package cli

import (
	"fmt"
	"io"
)

const usage = `orisan-review analyzes git diffs and patch files for review-routing risk.

Usage:
  orisan-review <command> [flags]

Commands:
  diff             Analyze a git diff
  scan-patch       Analyze a unified diff patch file or stdin
  list-rules       List rule catalogue
  list-categories  List sensitive change categories
  version          Print version

Flags:
  -h, --help       Show help
`

func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		_, _ = fmt.Fprint(stdout, usage)
		return ExitOK
	}

	switch args[0] {
	case "diff":
		return runDiff(args[1:], stdout, stderr)
	case "scan-patch":
		return runScanPatch(args[1:], stdin, stdout, stderr)
	case "list-rules":
		return runListRules(args[1:], stdout, stderr)
	case "list-categories":
		return runListCategories(args[1:], stdout, stderr)
	case "version":
		return runVersion(args[1:], stdout, stderr)
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		_, _ = fmt.Fprint(stderr, usage)
		return ExitUsageError
	}
}
