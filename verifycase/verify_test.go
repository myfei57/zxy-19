package verifycase

import (
	"os"
	"testing"

	"coldchain/internal/alert"
	"coldchain/internal/audit"
	"coldchain/internal/batch"
	"coldchain/internal/data"
	"coldchain/internal/flow"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestReleaseAfterFrozenClear verifies a batch is released only after its
// frozen flag is durably cleared.
func TestReleaseAfterFrozenClear(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	batches := batch.New(st)
	temps := temp.New(st)
	alerts := alert.New(st, temps, batches)
	audits := audit.New(st)
	flows := flow.New(st, batches, alerts, audits)
	b, err := batches.Register("B-08", "ns-01")
	if err != nil {
		t.Fatal(err)
	}
	if err := batches.MarkInStorage(b.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := alerts.Freeze(b.ID, data.TemperatureReading{ID: "r", Celsius: -15}); err != nil {
		t.Fatal(err)
	}
	flagPath := st.Path("frozen", b.ID+".json")
	if err := os.Chmod(flagPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if err := flows.Release(b.ID); err == nil {
		t.Fatal("release should fail when the frozen flag cannot be cleared")
	}
	got, err := batches.Get(b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status == data.StatusReleased {
		t.Fatalf("batch released before frozen flag durably cleared: status = %s", got.Status)
	}
}
