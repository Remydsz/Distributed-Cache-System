package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	httpsrv "distcache/internal/server/http"
	"distcache/internal/cache"
	_ "distcache/internal/metrics"
)

func main() {
	port := getenv("PORT", "8080")
	capacity := atoi(getenv("CAP", "200000"))

	lru := cache.NewLRU(capacity)
	srv := httpsrv.New(lru)

	log.Printf("cache-server listening on :%s (cap=%d)", port, capacity)
	log.Fatal(http.ListenAndServe(":"+port, srv.Routes()))
}

func getenv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }
func atoi(s string) int         { i, _ := strconv.Atoi(s); return i }
