package main

import (
	"time"

	"coldchain/internal/console"
	"coldchain/internal/data"
)

// seedDemo registers a demo namespace, probe and batch, then walks the full
// control chain so every durable ordering is exercised at least once.
func seedDemo(s *console.Services) error {
	nsv, err := s.NS.Register("北区冷藏库", "warehouse", "赵翔宇")
	if err != nil {
		return err
	}
	probe, err := s.Probe.Register(nsv.ID, "冷库-01", "")
	if err != nil {
		return err
	}
	batch, err := s.Batch.Register("B20260822-001", nsv.ID)
	if err != nil {
		return err
	}
	if err := s.Batch.MarkInStorage(batch.ID); err != nil {
		return err
	}
	if err := s.Quota.SetLimit(probe.ID, 500); err != nil {
		return err
	}
	now := time.Now().UTC()
	for i := 0; i < 30; i++ {
		reading := data.TemperatureReading{
			ID: data.NewID(), ProbeID: probe.ID, BatchID: batch.ID,
			Celsius: -20 + float64(i%5), RecordedAt: now.Add(time.Duration(i) * time.Minute),
		}
		if err := s.Probe.Report(probe.ID, reading); err != nil {
			return err
		}
	}
	if err := s.Probe.Send(probe.ID); err != nil {
		return err
	}
	win := data.Window{
		ID: data.NewID(), BatchID: batch.ID, ProbeID: probe.ID,
		Start: now, End: now.Add(30 * time.Minute), Generation: 1,
		Status: "open", ReadingCount: 30,
	}
	if _, err := s.Temp.Aggregate(batch.ID, win); err != nil {
		return err
	}
	if err := s.Temp.CloseWindow(batch.ID, win); err != nil {
		return err
	}
	rule, err := s.Rule.Publish("冷藏阈值", "celsius", -18)
	if err != nil {
		return err
	}
	verdict, err := s.Rule.Check(data.TemperatureReading{
		ID: data.NewID(), ProbeID: probe.ID, BatchID: batch.ID,
		Celsius: -15, RecordedAt: now.Add(31 * time.Minute),
	})
	if err != nil {
		return err
	}
	if verdict.Overheat {
		if _, err := s.Alert.Freeze(batch.ID, data.TemperatureReading{
			ID: data.NewID(), ProbeID: probe.ID, BatchID: batch.ID,
			Celsius: -15, RecordedAt: now.Add(31 * time.Minute),
		}); err != nil {
			return err
		}
	}
	if err := s.Flow.Release(batch.ID); err != nil {
		return err
	}
	if _, err := s.Flow.Ship(batch.ID); err != nil {
		return err
	}
	if err := s.Flow.Receive(batch.ID); err != nil {
		return err
	}
	if _, err := s.Trace.Replay(batch.ID, rule.Version); err != nil {
		return err
	}
	if _, err := s.Trace.Timeline(batch.ID); err != nil {
		return err
	}
	if _, err := s.Audit.ListBySubject(batch.ID); err != nil {
		return err
	}
	return nil
}
