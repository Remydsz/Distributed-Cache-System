# DistCache

A lightweight **distributed in-memory cache** written in Go.  
Features **LRU eviction**, **TTL expiration**, and a simple **HTTP API**.  
Scales horizontally across nodes using **consistent hashing + virtual nodes (vnodes)** with client-side routing.

---

## Features
- In-memory key-value cache with:
  - **LRU eviction policy** (O(1) get/set)
  - **TTL expiration**
- **HTTP API** for `GET`, `PUT`, `DELETE`, and health checks
- **Prometheus metrics** for hits, misses, evictions, TTL expirations, and latency
- **Consistent hashing with vnodes** for even key distribution across nodes
- **CLI router** for transparent multi-node access
- Benchmarked at **80K+ QPS per node** on cache hits (sub-ms latency)

---

### How to run

To run a single server:
```bash
go run ./cmd/cache-server

Run multiple:
PORT=8080 go run ./cmd/cache-server
PORT=8081 go run ./cmd/cache-server
PORT=8082 go run ./cmd/cache-server

Benchmarked between 30k (cache read, write, misses) and 80k (read only) QPS

