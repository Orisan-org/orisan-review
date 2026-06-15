package rules

import (
	"regexp"
	"sort"
	"strings"

	"github.com/orisan/review/internal/classify"
	"github.com/orisan/review/internal/model"
	"github.com/orisan/review/internal/redact"
)

type Engine struct{}

func (Engine) Run(files []model.ChangedFile) []model.Finding {
	var findings []model.Finding
	for _, file := range files {
		if file.IsBinary {
			findings = append(findings, finding("REVIEW-BIN-001", "Binary file change detected", model.SeverityMedium, "binary_file_change", file.NewPath, 0, "Binary diff cannot be inspected safely.", "Binary changes can hide behavior changes from text review.", "Review binary artifact provenance before merge.", "medium", ""))
			continue
		}
		findings = append(findings, runFileRules(file)...)
	}
	sortFindings(findings)
	return findings
}

func runFileRules(file model.ChangedFile) []model.Finding {
	var findings []model.Finding
	path := file.NewPath
	if path == "" {
		path = file.OldPath
	}

	if isDependencyManifest(path) {
		findings = append(findings, finding("REVIEW-DEP-070", "Dependency manifest changed", model.SeverityMedium, "dependency_manifest_changed", path, 0, "Dependency manifest or lockfile changed.", "Dependency set changed and should be reviewed before merge.", "Review package additions/removals. Do not assume vulnerability without advisory lookup.", "high", ""))
	}

	for _, line := range changedLines(file) {
		lower := strings.ToLower(line.Content)
		switch {
		case hasAIMarker(lower):
			findings = append(findings, finding("REVIEW-AI-001", "AI-assistance marker detected", model.SeverityMedium, "ai_generated_marker", path, lineNumber(line), "AI-assistance marker detected in changed content.", "Marker suggests generated or agent-assisted output may need closer validation.", "Validate the changed output before trusting it.", "medium", line.Content))
		case isAuthLogic(path, lower):
			findings = append(findings, finding("REVIEW-AUTH-010", "Authentication logic changed", model.SeverityHigh, "auth_logic_changed", path, lineNumber(line), "Authentication-sensitive logic changed.", "Authentication behavior changed and should receive security review.", "Review auth/session/token behavior before merge.", "medium", line.Content))
		}
	}

	if weakenedAuthorization(file) {
		line := firstRelevantLine(file, "return nil", "return true", "allow", "skip")
		findings = append(findings, finding("REVIEW-AUTH-011", "Authorization check removed or weakened", model.SeverityCritical, "authorization_weakened", path, lineNumber(line), "Authorization-related check appears removed or weakened.", "Access control may allow actions that were previously denied.", "Review the authorization condition and require explicit approval before merge.", "high", line.Content))
	}
	if validationRemoved(file) {
		line := firstRelevantLine(file, "validate", "sanitize", "schema.parse", "skip validation")
		findings = append(findings, finding("REVIEW-VAL-020", "Input validation removed or weakened", model.SeverityHigh, "validation_removed", path, lineNumber(line), "Input validation appears removed or bypassed.", "Unvalidated input may reach sensitive code paths.", "Restore validation or document a reviewed replacement.", "high", line.Content))
	}
	if tlsDisabled(file) {
		line := firstRelevantLine(file, "insecureskipverify", "rejectunauthorized", "verify=false", "curl -k")
		findings = append(findings, finding("REVIEW-CRYPTO-030", "TLS verification disabled", model.SeverityCritical, "tls_verification_disabled", path, lineNumber(line), "Added line disables TLS verification.", "Outbound HTTPS connections may accept invalid certificates.", "Do not disable TLS verification. Use proper CA configuration for test or private environments.", "high", line.Content))
	}
	if secretAdded(file) {
		line := firstSecretLine(file)
		findings = append(findings, finding("REVIEW-SEC-040", "Secret-like value added", model.SeverityCritical, "secret_like_value_added", path, lineNumber(line), "Added line contains a secret-like value.", "Secrets committed to code can be copied, logged, or abused.", "Remove the secret, rotate it if real, and use a secret manager.", "high", line.Content))
	}
	if commandExecutionAdded(file) {
		line := firstCommandExecutionLine(file)
		findings = append(findings, finding("REVIEW-RCE-050", "Command execution surface added", model.SeverityCritical, "command_execution_added", path, lineNumber(line), "Added line introduces command execution or a shell pipeline.", "Command execution surfaces can turn untrusted input or agent output into host code execution.", "Require security review and constrain executable inputs before merge.", "high", line.Content))
	}
	if ciPermissionsBroadened(file) {
		line := firstRelevantLine(file, "write-all", "contents: write", "id-token: write", "actions: write", "pull-requests: write")
		findings = append(findings, finding("REVIEW-CICD-060", "CI workflow permissions broadened", model.SeverityHigh, "ci_permissions_broadened", path, lineNumber(line), "GitHub Actions permissions changed toward broader write access.", "Broader workflow token permissions can expand CI/CD blast radius.", "Use least privilege and require CI/CD owner review.", "high", line.Content))
	}
	if unpinnedAction(file) {
		line := firstRelevantLine(file, "uses:")
		findings = append(findings, finding("REVIEW-CICD-061", "Unpinned third-party action introduced", model.SeverityMedium, "unpinned_github_action", path, lineNumber(line), "Third-party GitHub Action is not pinned to a full commit SHA.", "Mutable tags or branches can change after review.", "Pin third-party actions to a full-length commit SHA when practical.", "medium", line.Content))
	}
	if infraPublicExposure(file) {
		line := firstRelevantLine(file, "0.0.0.0/0", "::/0", "loadbalancer", "nodeport")
		findings = append(findings, finding("REVIEW-INFRA-080", "Public network exposure changed", model.SeverityHigh, "infra_public_exposure", path, lineNumber(line), "Infrastructure change appears to expose a service or network path publicly.", "Public exposure can increase attack surface.", "Confirm the exposure is intentional and reviewed by infra/security.", "high", line.Content))
	}
	if destructiveMigration(file) {
		line := firstRelevantLine(file, "drop table", "drop column", "truncate", "delete from")
		findings = append(findings, finding("REVIEW-DB-090", "Destructive migration", model.SeverityHigh, "destructive_migration", path, lineNumber(line), "Migration contains destructive database operation.", "Data may be removed or made unrecoverable.", "Require human review and rollback planning before merge.", "high", line.Content))
	}
	if testsSkipped(file) {
		line := firstRelevantLine(file, "test.skip", "pytest.mark.skip", "t.skip", "xdescribe", "xit")
		findings = append(findings, finding("REVIEW-TEST-100", "Tests removed or skipped", model.SeverityMedium, "tests_skipped", path, lineNumber(line), "Changed test code appears to skip or disable tests.", "Review gates may be weakened.", "Confirm skipped tests are intentional and temporary.", "medium", line.Content))
	}

	return dedupeFindings(findings)
}

func changedLines(file model.ChangedFile) []model.DiffLine {
	var lines []model.DiffLine
	for _, hunk := range file.Hunks {
		for _, line := range hunk.Lines {
			if line.Type == model.LineAdded || line.Type == model.LineRemoved {
				lines = append(lines, line)
			}
		}
	}
	return lines
}

func addedLines(file model.ChangedFile) []model.DiffLine {
	var lines []model.DiffLine
	for _, hunk := range file.Hunks {
		for _, line := range hunk.Lines {
			if line.Type == model.LineAdded {
				lines = append(lines, line)
			}
		}
	}
	return lines
}

func removedLines(file model.ChangedFile) []model.DiffLine {
	var lines []model.DiffLine
	for _, hunk := range file.Hunks {
		for _, line := range hunk.Lines {
			if line.Type == model.LineRemoved {
				lines = append(lines, line)
			}
		}
	}
	return lines
}

func hasAIMarker(lower string) bool {
	return containsAny(lower, "generated by chatgpt", "generated by claude", "generated with copilot", "ai-generated", "llm-generated", "this code was generated", "codex", "claude-code", "cursor agent", "windsurf")
}

func isAuthLogic(path, lower string) bool {
	if !hasCategory(path, "auth") {
		return false
	}
	return containsAny(lower, "login", "logout", "authenticate", "authorization", "authorize", "jwt", "token", "session", "cookie", "oauth", "password", "mfa", "2fa")
}

func weakenedAuthorization(file model.ChangedFile) bool {
	removed := strings.ToLower(joinLineContents(removedLines(file)))
	added := strings.ToLower(joinLineContents(addedLines(file)))
	return containsAny(removed, "requireadmin", "isadmin", "haspermission", "authorize", "checkpermission", "canaccess", "forbidden", "403") &&
		containsAny(added, "return nil", "return true", "allow", "skip", "todo")
}

func validationRemoved(file model.ChangedFile) bool {
	removed := strings.ToLower(joinLineContents(removedLines(file)))
	added := strings.ToLower(joinLineContents(addedLines(file)))
	return containsAny(removed, "validate", "sanitize", "escape", "schema.parse", "zod.parse", "joi.validate", "validator", "isvalid") ||
		containsAny(added, "todo validate later", "skip validation", "validate: false", "ignore validation")
}

func tlsDisabled(file model.ChangedFile) bool {
	added := strings.ToLower(joinLineContents(addedLines(file)))
	return containsAny(added, "insecureskipverify: true", "rejectunauthorized: false", "verify=false", "ssl._create_unverified_context", "node_tls_reject_unauthorized=0", "curl -k", "--insecure")
}

func secretAdded(file model.ChangedFile) bool {
	for _, line := range addedLines(file) {
		if secretPattern.MatchString(line.Content) {
			return true
		}
	}
	return false
}

func commandExecutionAdded(file model.ChangedFile) bool {
	return firstCommandExecutionLine(file).Content != ""
}

func firstCommandExecutionLine(file model.ChangedFile) model.DiffLine {
	for _, line := range addedLines(file) {
		if commandExecutionPattern.MatchString(line.Content) || curlPipeShellPattern.MatchString(line.Content) {
			return line
		}
	}
	return model.DiffLine{}
}

func ciPermissionsBroadened(file model.ChangedFile) bool {
	if !strings.HasPrefix(strings.ToLower(file.NewPath), ".github/workflows/") {
		return false
	}
	added := strings.ToLower(joinLineContents(addedLines(file)))
	return containsAny(added, "permissions: write-all", "write-all", "contents: write", "pull-requests: write", "actions: write", "id-token: write", "issues: write")
}

func unpinnedAction(file model.ChangedFile) bool {
	if !strings.HasPrefix(strings.ToLower(file.NewPath), ".github/workflows/") {
		return false
	}
	fullSHA := regexp.MustCompile(`@[a-f0-9]{40}\b`)
	for _, line := range addedLines(file) {
		lower := strings.ToLower(strings.TrimSpace(line.Content))
		lower = strings.TrimPrefix(lower, "- ")
		if strings.HasPrefix(lower, "uses: ") && !fullSHA.MatchString(lower) {
			return true
		}
	}
	return false
}

func isDependencyManifest(path string) bool {
	return hasCategory(path, "dependency_manifest")
}

func hasCategory(path, want string) bool {
	for _, category := range classify.CategoriesForPath(path) {
		if category == want {
			return true
		}
	}
	return false
}

func infraPublicExposure(file model.ChangedFile) bool {
	if !classify.IsSensitivePath(file.NewPath) {
		return false
	}
	added := strings.ToLower(joinLineContents(addedLines(file)))
	return containsAny(added, "0.0.0.0/0", "::/0", "nodeport", "loadbalancer", "ingress allow all")
}

func destructiveMigration(file model.ChangedFile) bool {
	added := strings.ToLower(joinLineContents(addedLines(file)))
	return containsAny(added, "drop table", "drop column", "truncate", "delete from", "alter table") && strings.Contains(added, "drop")
}

func testsSkipped(file model.ChangedFile) bool {
	added := strings.ToLower(joinLineContents(addedLines(file)))
	return containsAny(added, "test.skip", "pytest.mark.skip", "t.skip", "xdescribe", "xit")
}

func firstSecretLine(file model.ChangedFile) model.DiffLine {
	for _, line := range addedLines(file) {
		if secretPattern.MatchString(line.Content) {
			return line
		}
	}
	return model.DiffLine{}
}

func firstRelevantLine(file model.ChangedFile, needles ...string) model.DiffLine {
	for _, line := range changedLines(file) {
		lower := strings.ToLower(line.Content)
		for _, needle := range needles {
			if strings.Contains(lower, strings.ToLower(needle)) {
				return line
			}
		}
	}
	return model.DiffLine{}
}

func lineNumber(line model.DiffLine) int {
	if line.NewLine != 0 {
		return line.NewLine
	}
	return line.OldLine
}

func finding(id, title string, severity model.Severity, category, path string, line int, evidence, impact, remediation, confidence, snippet string) model.Finding {
	f := model.NewFinding(id, title, severity)
	f.Category = category
	f.Location = model.Location{Path: path, StartLine: line}
	f.Evidence = redact.Snippet(snippet, evidence)
	f.Impact = impact
	f.Remediation = remediation
	f.Confidence = confidence
	f.PayloadStored = false
	return f
}

func sortFindings(findings []model.Finding) {
	sort.Slice(findings, func(i, j int) bool {
		if model.SeverityRank(findings[i].Severity) != model.SeverityRank(findings[j].Severity) {
			return model.SeverityRank(findings[i].Severity) > model.SeverityRank(findings[j].Severity)
		}
		if findings[i].Location.Path != findings[j].Location.Path {
			return findings[i].Location.Path < findings[j].Location.Path
		}
		if findings[i].Location.StartLine != findings[j].Location.StartLine {
			return findings[i].Location.StartLine < findings[j].Location.StartLine
		}
		return findings[i].ID < findings[j].ID
	})
}

func dedupeFindings(findings []model.Finding) []model.Finding {
	seen := map[string]bool{}
	var out []model.Finding
	for _, finding := range findings {
		key := finding.ID + "|" + finding.Category + "|" + finding.Location.Path
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, finding)
	}
	return out
}

func joinLineContents(lines []model.DiffLine) string {
	var b strings.Builder
	for _, line := range lines {
		b.WriteString(line.Content)
		b.WriteByte('\n')
	}
	return b.String()
}

func containsAny(s string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(s, needle) {
			return true
		}
	}
	return false
}

var secretPattern = regexp.MustCompile(`(?i)(aws_secret_access_key|private_key|api[_-]?key|password|database_url|bearer\s+[a-z0-9._\-]+|gh[pousr]_[a-z0-9_]+|xox[baprs]-)`)
var commandExecutionPattern = regexp.MustCompile(`(?i)\b(os/exec|exec\.Command|subprocess\.(run|popen|call)|child_process\.(exec|spawn|execFile)|Runtime\.getRuntime\(\)\.exec|ProcessBuilder|system\(|popen\()`)
var curlPipeShellPattern = regexp.MustCompile(`(?i)\b(curl|wget)\b[^|\n\r]*\|\s*(sh|bash|zsh|fish|powershell|pwsh)\b`)
