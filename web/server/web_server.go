package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/roster"
	"github.com/fluxynet/goqa/subscriber/sse"
	"github.com/fluxynet/goqa/web"
)

type Server struct {
	Prefix    string
	Cache     goqa.Cache
	Broker    goqa.Broker
	Roster    goqa.Roster
	IndexHTML []byte
}

// List endpoint for coverage list api endpoint
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	var keys, err = s.Cache.Keys()
	if err != nil {
		web.JsonError(w, http.StatusInternalServerError, err)
		return
	} else if keys == nil {
		keys = []string{}
	}

	web.Json(w, keys)
}

// Get endpoint for single coverage api endpoint
func (s *Server) Get(w http.ResponseWriter, r *http.Request) {
	var pkg = strings.TrimPrefix(r.URL.Path, s.Prefix)

	var cov, ok = s.Cache.Get(pkg)

	if !ok {
		web.JsonError(w, http.StatusNotFound, web.ErrResourceNotFound)
		return
	}

	web.Json(w, cov)
}

// Index endpoint for proper display of IndexHTML index page
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	web.Print(w, http.StatusOK, web.ContentTypeHTML, s.IndexHTML)
}

// SSE endpoint for events updates
func (s *Server) SSE(w http.ResponseWriter, r *http.Request) {
	var flusher, ok = w.(http.Flusher)
	if !ok {
		web.JsonError(w, http.StatusPreconditionFailed, web.ErrStreamingNotSupported)
		return
	}

	w.Header().Set("Content-Type", web.ContentTypeEventStream)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	var (
		ctx = r.Context()
		sub = sse.New(w, flusher)
	)

	var err = s.Roster.Subscribe(ctx, goqa.EventCoverage, sub)
	if err != nil {
		fmt.Fprintf(w, "event: error")
		fmt.Fprintf(w, "data: failed to subscribe to event")
		return
	}

	roster.WatchCtx(ctx, s.Roster, sub)
}
