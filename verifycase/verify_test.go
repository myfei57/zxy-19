package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"coldchain/internal/alert"
	"coldchain/internal/batch"
	"coldchain/internal/data"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestFreezeAfterOverheatDurable verifies a batch freezes only after its
// overheat record is durable.
func TestFreezeAfterOverheatDurable(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	batches := batch.New(st)
	temps := temp.New(st)
	alerts := alert.New(st, temps, batches)
	b, err := batches.Register("B-02", "ns-01")
	if err != nil {
		t.Fatal(err)
	}
	if err := batches.MarkInStorage(b.ID); err != nil {
		t.Fatal(err)
	}
	overheatPath := st.Path("temp", "overheat", b.ID+".jsonl")
	if err := os.MkdirAll(filepath.Dir(overheatPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(overheatPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(overheatPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := alerts.Freeze(b.ID, data.TemperatureReading{ID: "r", Celsius: -15}); err == nil {
		t.Fatal("freeze should fail when the overheat record cannot be written")
	}
	got, err := batches.Get(b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status == data.StatusFrozen {
		t.Fatalf("batch frozen before overheat evidence durable: status = %s", got.Status)
	}
}
