package main

import (
	"coldchain/internal/alert"
	"coldchain/internal/audit"
	"coldchain/internal/batch"
	"coldchain/internal/console"
	"coldchain/internal/flow"
	"coldchain/internal/ns"
	"coldchain/internal/probe"
	"coldchain/internal/quota"
	"coldchain/internal/rule"
	"coldchain/internal/store"
	"coldchain/internal/temp"
	"coldchain/internal/trace"
)

// buildServices wires every domain service onto one file-backed store.
func buildServices(st *store.Store) *console.Services {
	nsSvc := ns.New(st)
	batches := batch.New(st)
	temps := temp.New(st)
	quotaSvc := quota.New(st)
	probes := probe.New(st, temps, quotaSvc)
	rules := rule.New(st, temps)
	alerts := alert.New(st, temps, batches)
	audits := audit.New(st)
	flows := flow.New(st, batches, alerts, audits)
	traces := trace.New(st, temps)
	return &console.Services{
		NS:    nsSvc,
		Probe: probes,
		Batch: batches,
		Temp:  temps,
		Rule:  rules,
		Alert: alerts,
		Flow:  flows,
		Trace: traces,
		Quota: quotaSvc,
		Audit: audits,
	}
}
