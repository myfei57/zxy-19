package console

import (
	"coldchain/internal/alert"
	"coldchain/internal/audit"
	"coldchain/internal/batch"
	"coldchain/internal/flow"
	"coldchain/internal/ns"
	"coldchain/internal/probe"
	"coldchain/internal/quota"
	"coldchain/internal/rule"
	"coldchain/internal/temp"
	"coldchain/internal/trace"
)

// Services bundles every domain service used by the console.
type Services struct {
	NS    *ns.Service
	Probe *probe.Service
	Batch *batch.Service
	Temp  *temp.Service
	Rule  *rule.Service
	Alert *alert.Service
	Flow  *flow.Service
	Trace *trace.Service
	Quota *quota.Service
	Audit *audit.Service
}
