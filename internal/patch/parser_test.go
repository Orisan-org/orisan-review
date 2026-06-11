package patch

import "testing"

func TestParseAddedRemovedContextLines(t *testing.T) {
	doc, err := Parse([]byte(`diff --git a/app/auth.go b/app/auth.go
index 1111111..2222222 100644
--- a/app/auth.go
+++ b/app/auth.go
@@ -10,3 +10,4 @@ func check() {
 keep()
-deny()
+allow()
+audit()
 }
`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(doc.Files) != 1 {
		t.Fatalf("files = %d, want 1", len(doc.Files))
	}
	file := doc.Files[0]
	if file.NewPath != "app/auth.go" || file.Status != "modified" {
		t.Fatalf("file = %+v", file)
	}
	if len(file.Hunks) != 1 {
		t.Fatalf("hunks = %d, want 1", len(file.Hunks))
	}
	lines := file.Hunks[0].Lines
	if len(lines) != 5 {
		t.Fatalf("lines = %d, want 5", len(lines))
	}
	if lines[1].OldLine != 11 || lines[1].Type != "removed" {
		t.Fatalf("removed line = %+v", lines[1])
	}
	if lines[2].NewLine != 11 || lines[2].Type != "added" {
		t.Fatalf("added line = %+v", lines[2])
	}
}

func TestParseRenamedDeletedAndBinaryFiles(t *testing.T) {
	doc, err := Parse([]byte(`diff --git a/old.go b/new.go
similarity index 90%
rename from old.go
rename to new.go
--- a/old.go
+++ b/new.go
@@ -1 +1 @@
-old
+new
diff --git a/dead.go b/dead.go
deleted file mode 100644
--- a/dead.go
+++ /dev/null
@@ -1 +0,0 @@
-dead
diff --git a/logo.png b/logo.png
Binary files a/logo.png and b/logo.png differ
`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(doc.Files) != 3 {
		t.Fatalf("files = %d, want 3", len(doc.Files))
	}
	if doc.Files[0].Status != "renamed" || doc.Files[0].OldPath != "old.go" || doc.Files[0].NewPath != "new.go" {
		t.Fatalf("renamed file = %+v", doc.Files[0])
	}
	if doc.Files[1].Status != "deleted" || doc.Files[1].NewPath != "dead.go" {
		t.Fatalf("deleted file = %+v", doc.Files[1])
	}
	if !doc.Files[2].IsBinary {
		t.Fatalf("binary file = %+v", doc.Files[2])
	}
}

func TestParseGitBinaryPatch(t *testing.T) {
	doc, err := Parse([]byte(`diff --git a/assets/logo.png b/assets/logo.png
index 8352675d67aed6625ece79af41c27fdb4ee2e867..eaf36c1daccfdf325514461cd1a2ffbc139b5464 100644
GIT binary patch
literal 4
LcmZQzWMT#Y01f~L

literal 3
KcmZQzWC8#H2LJ>B

`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(doc.Files) != 1 {
		t.Fatalf("files = %d, want 1", len(doc.Files))
	}
	if !doc.Files[0].IsBinary {
		t.Fatalf("git binary patch was not marked binary: %+v", doc.Files[0])
	}
	if doc.Files[0].NewPath != "assets/logo.png" {
		t.Fatalf("new path = %q, want assets/logo.png", doc.Files[0].NewPath)
	}
}

func TestParseGitBinaryPatchDelta(t *testing.T) {
	doc, err := Parse([]byte(`diff --git a/assets/logo.png b/assets/logo.png
index 2e65efe2a145dda7ee51d1741299f848e5bf752e..78981922613b2afb6025042ff6bd878ac1994e85 100644
GIT binary patch
delta 12
Tcma!nU|?VbVrqC?0Rsm?5L;Si

delta 8
PcmbQmU|?VbVrJ4a1ONa5

`))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(doc.Files) != 1 {
		t.Fatalf("files = %d, want 1", len(doc.Files))
	}
	if !doc.Files[0].IsBinary {
		t.Fatalf("git binary delta patch was not marked binary: %+v", doc.Files[0])
	}
}

func TestParseInvalidNonEmptyInput(t *testing.T) {
	_, err := Parse([]byte("not a patch\n"))
	if err != ErrNoUnifiedDiff {
		t.Fatalf("err = %v, want %v", err, ErrNoUnifiedDiff)
	}
}

func TestParseEmptyInput(t *testing.T) {
	doc, err := Parse(nil)
	if err != nil {
		t.Fatalf("Parse(nil) error = %v", err)
	}
	if len(doc.Files) != 0 {
		t.Fatalf("files = %d, want 0", len(doc.Files))
	}
}

func TestParseMalformedUnifiedDiff(t *testing.T) {
	_, err := Parse([]byte(`diff --git a/app.go b/app.go
index 1111111..2222222 100644
--- a/app.go
+++ b/app.go
@@ broken hunk
+fmt.Println("hello")
`))
	if err == nil {
		t.Fatal("expected malformed hunk error")
	}
}
