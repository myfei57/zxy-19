package console

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"coldchain/internal/data"
	"github.com/go-chi/chi/v5"
)

// writeError maps domain errors to HTTP status codes.
func (s *Server) writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, data.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, data.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, data.ErrQuotaExceeded):
		status = http.StatusTooManyRequests
	case errors.Is(err, data.ErrInvalidState), errors.Is(err, data.ErrFrozen):
		status = http.StatusConflict
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

type createNamespaceRequest struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Owner string `json:"owner"`
}

func (s *Server) createNamespace(w http.ResponseWriter, r *http.Request) {
	var req createNamespaceRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	nsv, err := s.ns.Register(req.Name, req.Kind, req.Owner)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, nsv)
}

func (s *Server) listNamespaces(w http.ResponseWriter, r *http.Request) {
	items, err := s.ns.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) deleteNamespace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.ns.Remove(id); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type createProbeRequest struct {
	NamespaceID string `json:"namespace_id"`
	Name        string `json:"name"`
	BatchID     string `json:"batch_id"`
}

func (s *Server) createProbe(w http.ResponseWriter, r *http.Request) {
	var req createProbeRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	p, err := s.probe.Register(req.NamespaceID, req.Name, req.BatchID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

func (s *Server) listProbes(w http.ResponseWriter, r *http.Request) {
	if namespaceID := r.URL.Query().Get("namespace_id"); namespaceID != "" {
		items, err := s.probe.ListByNamespace(namespaceID)
		if err != nil {
			s.writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	items, err := s.probe.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) probeCursor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	state, err := s.probe.Cursor(id)
	if err != nil {
		s.writeError(w, err)
		return
	}
	count, err := s.temp.ReadingCount(id)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"cursor": state, "readings": count})
}

func (s *Server) sendProbe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.probe.Send(id); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type reportRequest struct {
	Celsius    float64 `json:"celsius"`
	RecordedAt string  `json:"recorded_at"`
}

func (s *Server) reportReading(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req reportRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	reading := data.TemperatureReading{ID: data.NewID(), Celsius: req.Celsius}
	if req.RecordedAt != "" {
		recorded, err := time.Parse(time.RFC3339, req.RecordedAt)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		reading.RecordedAt = recorded
	} else {
		reading.RecordedAt = time.Now().UTC()
	}
	if err := s.probe.Report(id, reading); err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, reading)
}

type createBatchRequest struct {
	Code        string `json:"code"`
	NamespaceID string `json:"namespace_id"`
}

func (s *Server) createBatch(w http.ResponseWriter, r *http.Request) {
	var req createBatchRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	b, err := s.batch.Register(req.Code, req.NamespaceID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (s *Server) listBatches(w http.ResponseWriter, r *http.Request) {
	items, err := s.batch.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) getBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	b, err := s.batch.Get(id)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) batchNotices(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	items, err := s.flow.Notices(id)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) storeBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.batch.MarkInStorage(id); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) shipBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	notice, err := s.flow.Ship(id)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, notice)
}

func (s *Server) releaseBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.flow.Release(id); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) receiveBatch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.flow.Receive(id); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) listReadings(w http.ResponseWriter, r *http.Request) {
	probeID := r.URL.Query().Get("probe_id")
	if probeID == "" {
		items, err := s.temp.ListAllReadings()
		if err != nil {
			s.writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	items, err := s.temp.ListReadings(probeID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) batchCursor(w http.ResponseWriter, r *http.Request) {
	batchID := r.URL.Query().Get("batch_id")
	state, err := s.temp.Cursor(batchID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (s *Server) windowSummaries(w http.ResponseWriter, r *http.Request) {
	batchID := r.URL.Query().Get("batch_id")
	items, err := s.temp.Summaries(batchID)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type windowRequest struct {
	BatchID    string `json:"batch_id"`
	ProbeID    string `json:"probe_id"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Generation int    `json:"generation"`
	Status     string `json:"status"`
	Count      int    `json:"reading_count"`
}

func (s *Server) parseWindow(req windowRequest) (data.Window, error) {
	win := data.Window{
		ID:           data.NewID(),
		BatchID:      req.BatchID,
		ProbeID:      req.ProbeID,
		Generation:   req.Generation,
		Status:       req.Status,
		ReadingCount: req.Count,
	}
	if req.Start == "" || req.End == "" {
		return data.Window{}, data.ErrInvalidInput
	}
	start, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		return data.Window{}, err
	}
	end, err := time.Parse(time.RFC3339, req.End)
	if err != nil {
		return data.Window{}, err
	}
	win.Start = start
	win.End = end
	return win, nil
}

func (s *Server) aggregateWindow(w http.ResponseWriter, r *http.Request) {
	var req windowRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	win, err := s.parseWindow(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	result, err := s.temp.Aggregate(win.BatchID, win)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (s *Server) closeWindow(w http.ResponseWriter, r *http.Request) {
	var req windowRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	win, err := s.parseWindow(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.temp.CloseWindow(win.BatchID, win); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type publishRuleRequest struct {
	Name      string  `json:"name"`
	Metric    string  `json:"metric"`
	Threshold float64 `json:"threshold"`
}

func (s *Server) publishRule(w http.ResponseWriter, r *http.Request) {
	var req publishRuleRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	rule, err := s.rule.Publish(req.Name, req.Metric, req.Threshold)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, rule)
}

func (s *Server) getCurrentRule(w http.ResponseWriter, r *http.Request) {
	rule, err := s.rule.Current()
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rule)
}

type checkRequest struct {
	Celsius float64 `json:"celsius"`
}

func (s *Server) checkReading(w http.ResponseWriter, r *http.Request) {
	var req checkRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	verdict, err := s.rule.Check(data.TemperatureReading{ID: data.NewID(), Celsius: req.Celsius})
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, verdict)
}

type freezeRequest struct {
	BatchID string  `json:"batch_id"`
	Celsius float64 `json:"celsius"`
}

func (s *Server) freezeBatch(w http.ResponseWriter, r *http.Request) {
	var req freezeRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	alertRecord, err := s.alert.Freeze(req.BatchID, data.TemperatureReading{
		ID:      data.NewID(),
		Celsius: req.Celsius,
	})
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, alertRecord)
}

func (s *Server) listAlerts(w http.ResponseWriter, r *http.Request) {
	items, err := s.alert.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

type quotaRequest struct {
	Limit int `json:"limit"`
}

func (s *Server) getQuota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	state, err := s.quota.State(id)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, state)
}

func (s *Server) setQuota(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req quotaRequest
	if err := readJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.quota.SetLimit(id, req.Limit); err != nil {
		s.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) traceReplay(w http.ResponseWriter, r *http.Request) {
	batchID := r.URL.Query().Get("batch_id")
	if r.URL.Query().Get("generations") == "1" {
		count, err := s.temp.GenerationCount(batchID)
		if err != nil {
			s.writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"generations": count})
		return
	}
	if r.URL.Query().Get("timeline") == "1" {
		items, err := s.trace.Timeline(batchID)
		if err != nil {
			s.writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, items)
		return
	}
	generation := 0
	if raw := r.URL.Query().Get("generation"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		generation = parsed
	}
	items, err := s.trace.Replay(batchID, generation)
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) auditList(w http.ResponseWriter, r *http.Request) {
	items, err := s.audit.List()
	if err != nil {
		s.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
