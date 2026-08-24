package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) routes() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui/probes", http.StatusFound)
	})
	s.router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	s.router.Get("/ui/probes", s.page("probes.html"))
	s.router.Get("/ui/batches", s.page("batches.html"))
	s.router.Get("/ui/temperature", s.page("temperature.html"))
	s.router.Get("/ui/audit", s.page("audit.html"))

	s.router.Route("/api", func(r chi.Router) {
		r.Get("/status", s.status)
		r.Route("/namespaces", func(r chi.Router) {
			r.Get("/", s.listNamespaces)
			r.Post("/", s.createNamespace)
			r.Delete("/{id}", s.deleteNamespace)
		})
		r.Route("/probes", func(r chi.Router) {
			r.Get("/", s.listProbes)
			r.Post("/", s.createProbe)
			r.Get("/{id}/cursor", s.probeCursor)
			r.Post("/{id}/report", s.reportReading)
			r.Post("/{id}/send", s.sendProbe)
		})
		r.Route("/batches", func(r chi.Router) {
			r.Get("/", s.listBatches)
			r.Post("/", s.createBatch)
			r.Get("/{id}", s.getBatch)
			r.Get("/{id}/notices", s.batchNotices)
			r.Post("/{id}/store", s.storeBatch)
			r.Post("/{id}/ship", s.shipBatch)
			r.Post("/{id}/release", s.releaseBatch)
			r.Post("/{id}/receive", s.receiveBatch)
		})
		r.Route("/temperature", func(r chi.Router) {
			r.Get("/", s.listReadings)
			r.Get("/cursor", s.batchCursor)
			r.Get("/summaries", s.windowSummaries)
			r.Post("/aggregate", s.aggregateWindow)
			r.Post("/close", s.closeWindow)
		})
		r.Route("/rules", func(r chi.Router) {
			r.Get("/", s.getCurrentRule)
			r.Post("/", s.publishRule)
			r.Post("/check", s.checkReading)
		})
		r.Route("/alerts", func(r chi.Router) {
			r.Get("/", s.listAlerts)
			r.Post("/freeze", s.freezeBatch)
		})
		r.Route("/quota", func(r chi.Router) {
			r.Get("/{id}", s.getQuota)
			r.Post("/{id}", s.setQuota)
		})
		r.Route("/trace", func(r chi.Router) {
			r.Get("/", s.traceReplay)
		})
		r.Get("/audit", s.auditList)
	})
}
