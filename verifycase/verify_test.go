package verifycase

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"coldchain/internal/data"
	"coldchain/internal/probe"
	"coldchain/internal/quota"
	"coldchain/internal/store"
	"coldchain/internal/temp"
)

// TestProbeCursorAfterAckDurable verifies the probe cursor only advances
// after the sink acknowledgement is durable.
func TestProbeCursorAfterAckDurable(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	temps := temp.New(st)
	quotas := quota.New(st)
	probes := probe.New(st, temps, quotas)
	probeID := "probe-07"
	for i := 0; i < 3; i++ {
		if err := temps.AppendReading(probeID, data.TemperatureReading{
			ID: fmt.Sprintf("r%d", i), ProbeID: probeID, Celsius: -18,
		}); err != nil {
			t.Fatal(err)
		}
	}
	ackPath := st.Path("temp", "acks", probeID+".jsonl")
	if err := os.MkdirAll(filepath.Dir(ackPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ackPath, []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(ackPath, 0o444); err != nil {
		t.Fatal(err)
	}
	if err := probes.Send(probeID); err == nil {
		t.Fatal("send should fail when the ack cannot be written")
	}
	cursor, err := probes.Cursor(probeID)
	if err != nil {
		t.Fatal(err)
	}
	if cursor.Position != 0 {
		t.Fatalf("probe cursor advanced before ack durable: position = %d", cursor.Position)
	}
}
