package profile

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/pprof"
	"strconv"
)

// Handler returns an HTTP handler that serves the named profile.
func Handler(name string) http.Handler {
	// for goroutine profiler we should use debug flag = 1
	// because for flag = 2 it won't work
	if name == "goroutine" {
		return &handler{name: name, debug: 1}
	}
	return &handler{name: name, debug: 2}
}

// Handler - is a struct to create profile for the
type handler struct {
	name  string
	debug int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	p := pprof.Lookup(string(h.name))
	if p == nil {
		h.serveError(w, http.StatusNotFound, "Unknown profile")
		return
	}

	gc, _ := strconv.Atoi(r.FormValue("gc"))
	if h.name == "heap" && gc > 0 {
		runtime.GC()
	}

	p.WriteTo(w, h.debug)
}

func (h *handler) serveError(w http.ResponseWriter, status int, txt string) {
	w.Header().Set("X-Go-Pprof", "1")
	w.Header().Del("Content-Disposition")
	w.WriteHeader(status)
	fmt.Fprintln(w, txt)
}
