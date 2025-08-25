package httpsrv

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"distcache/internal/cache"
)

type Server struct {
	C *cache.LRU
}

func New(c *cache.LRU) *Server { return &Server{C: c} }

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200) })
	mux.HandleFunc("/cache/", s.handleCache)
	return mux
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
