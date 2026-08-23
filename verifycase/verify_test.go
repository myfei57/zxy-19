package verifycase

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestTempWindowAfterSummaryDurable verifies a window is marked closed only
// after its summary is durable.
func TestTempWindowAfterSummaryDurable(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	temps := temp.New(st)
	start := time.Now().Add(-time.Hour).UTC()
	win := data.Window{
		ID: "win-09", BatchID: "batch-09", ProbeID: "p1",
		Start: start, End: start.Add(time.Hour), Generation: 1,
		Status: "open", ReadingCount: 5,
	}
	summaryPath := st.Path("temp", "summaries", "batch-09.jsonl")
	if err := os.MkdirAll(filepath.Dir(summaryPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(summaryPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(summaryPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if err := temps.CloseWindow("batch-09", win); err == nil {
		t.Fatal("close should fail when the summary cannot be written")
	}
	got, err := temps.LoadWindow("batch-09", "win-09")
	if err == nil {
		if got.Status == temp.WindowClosed {
			t.Fatal("window closed before summary durable")
		}
		return
	}
	if !errors.Is(err, data.ErrNotFound) {
		t.Fatal(err)
	}
}
