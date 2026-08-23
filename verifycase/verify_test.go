package verifycase

import (
	"os"
	"testing"

	"coldchain/internal/alert"
	"coldchain/internal/audit"
	"coldchain/internal/batch"
	"coldchain/internal/flow"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestAuditAfterFlowDurable verifies the audit success event is recorded only
// after the flow state is durable.
func TestAuditAfterFlowDurable(t *testing.T) {
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
	b, err := batches.Register("B-10", "ns-01")
	if err != nil {
		t.Fatal(err)
	}
	if err := batches.MarkInStorage(b.ID); err != nil {
		t.Fatal(err)
	}
	batchPath := st.Path("batches", b.ID+".json")
	if err := os.Chmod(batchPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := flows.Ship(b.ID); err == nil {
		t.Fatal("ship should fail when the batch state cannot be written")
	}
	events, err := audits.ListBySubject(b.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, event := range events {
		if event.Actor == "ship" && event.OK {
			t.Fatal("audit success recorded before flow state durable")
		}
	}
}
