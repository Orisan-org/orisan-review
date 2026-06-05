package classify

import "strings"

func IsGeneratedPath(path string) bool {
	lower := strings.ToLower(path)
	return strings.Contains(lower, ".generated.") ||
		strings.Contains(lower, "/generated/") ||
		strings.Contains(lower, "/dist/") ||
		strings.Contains(lower, "/build/")
}
