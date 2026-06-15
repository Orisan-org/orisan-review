package classify

import (
	"path/filepath"
	"strings"
)

func IsSensitivePath(path string) bool {
	return len(CategoriesForPath(path)) > 0
}

func IsSecuritySensitivePath(path string) bool {
	for _, category := range CategoriesForPath(path) {
		if category != "public_claims" {
			return true
		}
	}
	return false
}

func CategoriesForPath(path string) []string {
	clean := filepath.ToSlash(strings.ToLower(path))
	base := filepath.Base(clean)
	var categories []string

	add := func(category string) {
		for _, existing := range categories {
			if existing == category {
				return
			}
		}
		categories = append(categories, category)
	}

	if containsAny(clean, "auth", "authentication", "authorization", "middleware", "session", "jwt", "oauth", "saml", "oidc", "rbac", "permission", "policy") {
		add("auth")
	}
	if containsAny(clean, "crypto", "tls", "cert") {
		add("crypto")
	}
	if containsAny(clean, "secret", ".env") {
		add("secrets")
	}
	if strings.HasPrefix(clean, ".github/workflows/") || base == ".gitlab-ci.yml" || base == "jenkinsfile" || base == "azure-pipelines.yml" || clean == ".circleci/config.yml" {
		add("ci_cd")
	}
	if isDependencyManifest(base) {
		add("dependency_manifest")
	}
	if isInfraPath(clean, base) {
		add("infrastructure")
	}
	if strings.HasPrefix(clean, "migrations/") || strings.HasPrefix(clean, "db/migrate/") || base == "schema.sql" || clean == "prisma/schema.prisma" {
		add("database_migration")
	}
	if strings.HasPrefix(clean, "docs/security/") || strings.HasPrefix(clean, "docs/privacy/") || strings.HasPrefix(clean, "website/") || base == "readme.md" || base == "security.md" || base == "privacy.md" || base == "terms.md" {
		add("public_claims")
	}
	if IsGeneratedPath(clean) {
		add("generated_code")
	}

	return categories
}

func containsAny(s string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(s, needle) {
			return true
		}
	}
	return false
}

func isDependencyManifest(base string) bool {
	switch base {
	case "package.json", "package-lock.json", "pnpm-lock.yaml", "yarn.lock", "requirements.txt", "pipfile", "pipfile.lock", "poetry.lock", "pyproject.toml", "go.mod", "go.sum", "cargo.toml", "cargo.lock", "pom.xml", "build.gradle", "composer.json", "gemfile", "gemfile.lock":
		return true
	default:
		return false
	}
}

func isInfraPath(path, base string) bool {
	return strings.HasSuffix(path, ".tf") ||
		strings.HasSuffix(path, ".tfvars") ||
		strings.HasPrefix(path, "kubernetes/") ||
		strings.HasPrefix(path, "k8s/") ||
		strings.HasPrefix(path, "helm/") ||
		strings.HasPrefix(path, "charts/") ||
		base == "dockerfile" ||
		base == "docker-compose.yml" ||
		base == "serverless.yml" ||
		strings.HasPrefix(path, "cloudformation/")
}
