package httpsrv

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"distcache/internal/cache"
	"distcache/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http/pprof"
)

type Server struct {
	C *cache.LRU
}

func New(c *cache.LRU) *Server { return &Server{C: c} }

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })

	// prometheus metrics
	mux.Handle("/metrics", promhttp.Handler())

	// pprof (optional but useful)
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// cache API with latency measurement
	mux.Handle("/cache/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		route := "/cache/"
		rec := &statusRecorder{ResponseWriter: w, code: 200}
		s.handleCache(rec, r)
		metrics.ObserveHTTP(r.Method, route, strconv.Itoa(rec.code), start)
	}))

	return mux
}

type statusRecorder struct {
	http.ResponseWriter
	code int
}
func (sr *statusRecorder) WriteHeader(code int) {
	sr.code = code
	sr.ResponseWriter.WriteHeader(code)
}


func (s *Server) handleCache(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/cache/"):]
	switch r.Method {
	case http.MethodGet:
		if v, ok := s.C.Get(key); ok {
			w.WriteHeader(200)
			_, _ = w.Write(v)
			return
		}
		http.NotFound(w, r)
	case http.MethodPut:
		body, _ := io.ReadAll(r.Body)
		ttl := parseTTL(r.URL.Query().Get("ttl"))
		s.C.Set(key, body, ttl)
		w.WriteHeader(204)
	case http.MethodDelete:
		s.C.Delete(key)
		w.WriteHeader(204)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func parseTTL(s string) time.Duration {
	if s == "" {
		return 0
	}
	// support "5s" or milliseconds as integer
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	if ms, err := strconv.Atoi(s); err == nil {
		return time.Duration(ms) * time.Millisecond
	}
	return 0
}
