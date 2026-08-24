package verifycase

import (
	"testing"

	"coldchain/internal/data"
	"coldchain/internal/rule"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestOverheatUsesCurrentRule verifies checks evaluate against the currently
// published rule version, never a pre-publication snapshot.
func TestOverheatUsesCurrentRule(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	temps := temp.New(st)
	rules := rule.New(st, temps)
	if _, err := rules.Publish("冷藏阈值", "celsius", -12); err != nil {
		t.Fatal(err)
	}
	if verdict, err := rules.Check(data.TemperatureReading{ID: "r1", Celsius: -8}); err != nil {
		t.Fatal(err)
	} else if !verdict.Overheat {
		t.Fatal("first check should overheat at -8 against -12")
	}
	if _, err := rules.Publish("冷藏阈值", "celsius", -18); err != nil {
		t.Fatal(err)
	}
	verdict, err := rules.Check(data.TemperatureReading{ID: "r2", Celsius: -15})
	if err != nil {
		t.Fatal(err)
	}
	if !verdict.Overheat {
		t.Fatalf("new threshold not used: -15 judged with version %d threshold %f", verdict.Version, verdict.Threshold)
	}
}
