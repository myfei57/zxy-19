package console

import (
	"net/http"
	"time"

	"coldchain/internal/alert"
	"coldchain/internal/audit"
	"coldchain/internal/batch"
	"coldchain/internal/config"
	"coldchain/internal/flow"
	"coldchain/internal/ns"
	"coldchain/internal/probe"
	"coldchain/internal/quota"
	"coldchain/internal/rule"
	"coldchain/internal/store"
	"coldchain/internal/temp"
	"coldchain/internal/trace"
	"github.com/go-chi/chi/v5"
)

// Server exposes the cold-chain console and JSON API.
type Server struct {
	opts   config.Options
	store  *store.Store
	ns     *ns.Service
	probe  *probe.Service
	batch  *batch.Service
	temp   *temp.Service
	rule   *rule.Service
	alert  *alert.Service
	flow   *flow.Service
	trace  *trace.Service
	quota  *quota.Service
	audit  *audit.Service
	router *chi.Mux
	http   *http.Server
}

// New builds a console server wired to every domain service.
func New(opts config.Options, st *store.Store, svcs *Services) *Server {
	s := &Server{
		opts:  opts,
		store: st,
		ns:    svcs.NS,
		probe: svcs.Probe,
		batch: svcs.Batch,
		temp:  svcs.Temp,
		rule:  svcs.Rule,
		alert: svcs.Alert,
		flow:  svcs.Flow,
		trace: svcs.Trace,
		quota: svcs.Quota,
		audit: svcs.Audit,
	}
	s.router = chi.NewRouter()
	s.routes()
	s.http = &http.Server{
		Addr:              opts.Addr,
		Handler:           s.router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return s
}

// Start serves HTTP until the process is interrupted.
func (s *Server) Start() error {
	return s.http.ListenAndServe()
}
