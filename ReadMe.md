# DistCache (v1)

Single-node in-memory cache with LRU + TTL and HTTP API.

## Run
```bash
go run ./cmd/cache-server
# Put with TTL=5s
curl -X PUT "localhost:8080/cache/foo?ttl=5s" -d 'hello'
# Get
curl localhost:8080/cache/foo
# Delete
curl -X DELETE localhost:8080/cache/foo
