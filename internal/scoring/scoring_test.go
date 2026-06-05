package scoring

import "testing"

func TestGradeScaffold(t *testing.T) {
	if Grade(nil) != "A" {
		t.Fatal("empty findings should grade A in scaffold")
	}
}
