package verifycase

import (
	"os"
	"path/filepath"
	"testing"

	"coldchain/internal/alert"
	"coldchain/internal/audit"
	"coldchain/internal/batch"
	"coldchain/internal/data"
	"coldchain/internal/flow"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestShipAfterDispatchDurable verifies a batch is marked shipped only after
// its dispatch notice is durable.
func TestShipAfterDispatchDurable(t *testing.T) {
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
	b, err := batches.Register("B-03", "ns-01")
	if err != nil {
		t.Fatal(err)
	}
	if err := batches.MarkInStorage(b.ID); err != nil {
		t.Fatal(err)
	}
	noticePath := st.Path("flow", "notices", b.ID+".jsonl")
	if err := os.MkdirAll(filepath.Dir(noticePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(noticePath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(noticePath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := flows.Ship(b.ID); err == nil {
		t.Fatal("ship should fail when the dispatch notice cannot be written")
	}
	got, err := batches.Get(b.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status == data.StatusShipped {
		t.Fatalf("batch shipped before dispatch notice durable: status = %s", got.Status)
	}
}
