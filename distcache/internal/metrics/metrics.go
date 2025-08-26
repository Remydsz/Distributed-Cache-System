package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	CacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "distcache_hits_total", Help: "Total cache hits",
	})
	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "distcache_misses_total", Help: "Total cache misses",
	})
	CacheEvictions = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "distcache_evictions_total", Help: "Total evictions from LRU",
	})
	CacheExpired = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "distcache_expired_total", Help: "Total keys expired due to TTL",
	})
	CacheSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "distcache_items", Help: "Current number of keys in cache",
	})
	HTTPDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "distcache_http_latency_seconds",
		Help:    "HTTP handler latency",
		Buckets: prometheus.DefBuckets, // ~[5ms..10s]
	}, []string{"method", "route", "status"})
)

func init() {
	prometheus.MustRegister(
		CacheHits, CacheMisses, CacheEvictions, CacheExpired, CacheSize, HTTPDuration,
	)
}

func ObserveHTTP(method, route, status string, start time.Time) {
	HTTPDuration.WithLabelValues(method, route, status).
		Observe(time.Since(start).Seconds())
}
