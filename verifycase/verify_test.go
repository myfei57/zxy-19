package verifycase

import (
	"testing"
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
	"coldchain/internal/temp"
	"coldchain/internal/trace"
)

// TestTraceReplaySingleGeneration verifies replay loads one window
// generation, never a cross-generation mix.
func TestTraceReplaySingleGeneration(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	temps := temp.New(st)
	traces := trace.New(st, temps)
	start := time.Now().Add(-time.Hour).UTC()
	old := data.Window{
		ID: "w1", BatchID: "batch-06", ProbeID: "p1",
		Start: start, End: start.Add(time.Hour), Generation: 1, Status: "closed",
	}
	current := data.Window{
		ID: "w2", BatchID: "batch-06", ProbeID: "p1",
		Start: start, End: start.Add(time.Hour), Generation: 2, Status: "closed",
	}
	if err := temps.SaveWindow(old); err != nil {
		t.Fatal(err)
	}
	if err := temps.SaveWindow(current); err != nil {
		t.Fatal(err)
	}
	windows, err := traces.Replay("batch-06", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(windows) != 1 || windows[0].Generation != 2 {
		t.Fatalf("replay mixed generations: %+v", windows)
	}
}
