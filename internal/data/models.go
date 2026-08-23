package data

import "time"

// Namespace identifies a warehouse, vehicle or container that owns probes.
type Namespace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
}

// Probe is a temperature sensor attached to a namespace and optional batch.
type Probe struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	NamespaceID string    `json:"namespace_id"`
	BatchID     string    `json:"batch_id"`
	Sequence    int       `json:"sequence"`
	CreatedAt   time.Time `json:"created_at"`
}

// Batch tracks one cold-chain lot through the lifecycle state machine.
type Batch struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	NamespaceID string    `json:"namespace_id"`
	Status      string    `json:"status"`
	Frozen      bool      `json:"frozen"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TemperatureReading is one probe sample persisted for aggregation.
type TemperatureReading struct {
	ID         string    `json:"id"`
	ProbeID    string    `json:"probe_id"`
	BatchID    string    `json:"batch_id"`
	Celsius    float64   `json:"celsius"`
	RecordedAt time.Time `json:"recorded_at"`
}

// OverheatRecord is the durable evidence written before a batch freezes.
type OverheatRecord struct {
	ID          string             `json:"id"`
	BatchID     string             `json:"batch_id"`
	Reading     TemperatureReading `json:"reading"`
	RuleVersion string             `json:"rule_version"`
	Threshold   float64            `json:"threshold"`
	CreatedAt   time.Time          `json:"created_at"`
}

// Window is one aggregation time slot of a batch generation.
type Window struct {
	ID           string    `json:"id"`
	BatchID      string    `json:"batch_id"`
	ProbeID      string    `json:"probe_id"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	Generation   int       `json:"generation"`
	Status       string    `json:"status"`
	ReadingCount int       `json:"reading_count"`
}

// WindowResult is the durable outcome of aggregating one window.
type WindowResult struct {
	ID         string    `json:"id"`
	BatchID    string    `json:"batch_id"`
	ProbeID    string    `json:"probe_id"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Generation int       `json:"generation"`
	Min        float64   `json:"min"`
	Max        float64   `json:"max"`
	Avg        float64   `json:"avg"`
	Count      int       `json:"count"`
	ComputedAt time.Time `json:"computed_at"`
}

// WindowSummary is the durable closing summary of a window.
type WindowSummary struct {
	ID           string    `json:"id"`
	BatchID      string    `json:"batch_id"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	Generation   int       `json:"generation"`
	Records      int       `json:"records"`
	OK           bool      `json:"ok"`
	SummarizedAt time.Time `json:"summarized_at"`
}

// Rule is one published overheat threshold version.
type Rule struct {
	ID          string    `json:"id"`
	Version     int       `json:"version"`
	Name        string    `json:"name"`
	Metric      string    `json:"metric"`
	Threshold   float64   `json:"threshold"`
	PublishedAt time.Time `json:"published_at"`
}

// Verdict is the outcome of checking one reading against a rule version.
type Verdict struct {
	Overheat  bool    `json:"overheat"`
	Reading   float64 `json:"reading"`
	Threshold float64 `json:"threshold"`
	Version   int     `json:"version"`
	Reason    string  `json:"reason,omitempty"`
}

// Alert records a freeze decision for one batch.
type Alert struct {
	ID        string    `json:"id"`
	BatchID   string    `json:"batch_id"`
	Severity  string    `json:"severity"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// Notice is a dispatch or arrival message durably sent to the receiver.
type Notice struct {
	ID        string    `json:"id"`
	BatchID   string    `json:"batch_id"`
	Kind      string    `json:"kind"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditEvent is an immutable operation record.
type AuditEvent struct {
	ID        string    `json:"id"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Subject   string    `json:"subject"`
	OK        bool      `json:"ok"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

// QuotaState is the durable per-probe reporting quota of one window.
type QuotaState struct {
	ProbeID   string    `json:"probe_id"`
	Limit     int       `json:"limit"`
	Used      int       `json:"used"`
	Window    string    `json:"window"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CursorState is the durable position of an aggregation or forwarding cursor.
type CursorState struct {
	ID        string    `json:"id"`
	Position  int       `json:"position"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SinkBatch is the acknowledgement a downstream sink durably accepts.
type SinkBatch struct {
	ID        string    `json:"id"`
	ProbeID   string    `json:"probe_id"`
	From      int       `json:"from"`
	To        int       `json:"to"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

// FrozenFlag is the durable marker that a batch is or is not frozen.
type FrozenFlag struct {
	BatchID   string    `json:"batch_id"`
	Frozen    bool      `json:"frozen"`
	UpdatedAt time.Time `json:"updated_at"`
}
