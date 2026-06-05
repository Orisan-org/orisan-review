package rules

import "testing"

func TestCatalogue(t *testing.T) {
	if len(Catalogue()) == 0 {
		t.Fatal("expected placeholder rule catalogue")
	}
}
