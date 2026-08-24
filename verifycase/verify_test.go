package verifycase

import (
	"testing"

	"coldchain/internal/data"
	"coldchain/internal/probe"
	"coldchain/internal/quota"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestProbeQuotaRejectsBeforeReport verifies an over-quota reading is never
// stored: the quota gate runs before the durable write.
func TestProbeQuotaRejectsBeforeReport(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	temps := temp.New(st)
	quotas := quota.New(st)
	probes := probe.New(st, temps, quotas)
	probeID := "probe-05"
	if err := quotas.SetLimit(probeID, 1); err != nil {
		t.Fatal(err)
	}
	if err := probes.Append(probeID, data.TemperatureReading{ID: "r1", Celsius: -20}); err != nil {
		t.Fatal(err)
	}
	count, err := temps.ReadingCount(probeID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("first reading not stored: count = %d", count)
	}
	if err := probes.Append(probeID, data.TemperatureReading{ID: "r2", Celsius: -21}); err == nil {
		t.Fatal("over-quota report should be rejected")
	}
	count, err = temps.ReadingCount(probeID)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("over-quota reading consumed storage: count = %d", count)
	}
}
