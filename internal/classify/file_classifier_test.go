package classify

import "testing"

func TestCategories(t *testing.T) {
	categories := Categories()
	if len(categories) == 0 {
		t.Fatal("expected scaffold categories")
	}
}

func TestCategoriesForPath(t *testing.T) {
	cases := map[string]string{
		"src/auth/session.go":           "auth",
		".github/workflows/deploy.yml":  "ci_cd",
		"package.json":                  "dependency_manifest",
		"infra/main.tf":                 "infrastructure",
		"migrations/001_drop_table.sql": "database_migration",
		"README.md":                     "public_claims",
	}
	for path, want := range cases {
		got := CategoriesForPath(path)
		found := false
		for _, category := range got {
			if category == want {
				found = true
			}
		}
		if !found {
			t.Fatalf("CategoriesForPath(%q) = %v, missing %q", path, got, want)
		}
	}
}
