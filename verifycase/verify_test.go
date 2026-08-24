package verifycase

import (
	"os"
	"testing"
	"time"

	"coldchain/internal/data"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestTempWindowAfterResultDurable verifies the window cursor only advances
// after the aggregated window result is durable.
func TestTempWindowAfterResultDurable(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	temps := temp.New(st)
	probeID := "probe-01"
	start := time.Now().Add(-time.Hour).UTC()
	win := data.Window{
		ID: "win-01", BatchID: "batch-01", ProbeID: probeID,
		Start: start, End: start.Add(time.Hour), Generation: 1, Status: "open",
	}
	if err := temps.AppendReading(probeID, data.TemperatureReading{
		ID: "r1", ProbeID: probeID, BatchID: "batch-01",
		Celsius: -20, RecordedAt: start.Add(time.Minute),
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := temps.Aggregate("batch-01", win); err != nil {
		t.Fatal(err)
	}
	resultPath := st.Path("temp", "results", "batch-01.jsonl")
	if err := os.Chmod(resultPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := temps.Aggregate("batch-01", win); err == nil {
		t.Fatal("aggregate should fail when the result cannot be written")
	}
	cursor, err := temps.Cursor("batch-01")
	if err != nil {
		t.Fatal(err)
	}
	if cursor.Position != 1 {
		t.Fatalf("cursor advanced before result durable: position = %d", cursor.Position)
	}
}
