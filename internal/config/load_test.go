package config

import "testing"

func TestDefault(t *testing.T) {
	if Default().SeverityThreshold != "high" {
		t.Fatal("default severity threshold should be high")
	}
}
