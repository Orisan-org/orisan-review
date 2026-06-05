package model

import "testing"

func TestSeverityRank(t *testing.T) {
	if SeverityRank(SeverityCritical) <= SeverityRank(SeverityHigh) {
		t.Fatal("critical should rank above high")
	}
	if SeverityRank(Severity("unknown")) != -1 {
		t.Fatal("unknown severity should rank -1")
	}
}

func TestNewFindingPayloadStoredDefault(t *testing.T) {
	finding := NewFinding("REVIEW-TEST-000", "test", SeverityLow)
	if finding.PayloadStored {
		t.Fatal("new finding must default payload_stored=false")
	}
}
