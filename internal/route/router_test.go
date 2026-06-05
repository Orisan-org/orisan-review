package route

import "testing"

func TestRoutesForFindingsScaffold(t *testing.T) {
	if routes := RoutesForFindings(nil); routes != nil {
		t.Fatalf("routes = %v, want nil", routes)
	}
}
