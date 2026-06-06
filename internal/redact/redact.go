package redact

import (
	"regexp"
	"strings"
)

const maxSnippetLength = 160

func Evidence(input string) string {
	return Snippet(input, "")
}

func Snippet(input, fallback string) string {
	snippet := strings.TrimSpace(input)
	if snippet == "" {
		snippet = fallback
	}
	snippet = redactSecrets(snippet)
	if len(snippet) > maxSnippetLength {
		snippet = snippet[:maxSnippetLength] + "..."
	}
	return snippet
}

func redactSecrets(input string) string {
	out := input
	for _, pattern := range secretPatterns {
		out = pattern.ReplaceAllString(out, `${1}REDACTED`)
	}
	out = privateKeyPattern.ReplaceAllString(out, "PRIVATE_KEY=REDACTED")
	return out
}

var secretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(AWS_SECRET_ACCESS_KEY\s*=\s*["']?)[^"'\s]+`),
	regexp.MustCompile(`(?i)(DATABASE_URL\s*=\s*["']?)[^"'\s]+`),
	regexp.MustCompile(`(?i)(api[_-]?key\s*[:=]\s*["']?)[^"'\s]+`),
	regexp.MustCompile(`(?i)(password\s*[:=]\s*["']?)[^"'\s]+`),
	regexp.MustCompile(`(?i)(bearer\s+)[A-Za-z0-9._\-]+`),
	regexp.MustCompile(`(?i)(gh[pousr]_)[A-Za-z0-9_]+`),
}

var privateKeyPattern = regexp.MustCompile(`(?i)PRIVATE_KEY\s*=\s*["']?-+BEGIN PRIVATE KEY-+.*`)
