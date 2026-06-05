package textutil

import "testing"

func TestSnippet(t *testing.T) {
	if got := Snippet("abcdef", 3); got != "abc" {
		t.Fatalf("Snippet() = %q, want abc", got)
	}
}
