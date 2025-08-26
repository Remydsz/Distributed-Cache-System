package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"distcache/internal/ring"
)

func main() {
	// Cluster nodes, comma-separated
	peersEnv := getenv("PEERS", "http://localhost:8080,http://localhost:8081,http://localhost:8082")
	peers := strings.Split(peersEnv, ",")

	nodes := make([]ring.NodeID, 0, len(peers))
	for _, p := range peers {
		nodes = append(nodes, ring.NodeID(strings.TrimSpace(p)))
	}

	// 3 vnodes per node
	r := ring.New(nodes, 3)

	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "set":
		fs := flag.NewFlagSet("set", flag.ExitOnError)
		ttl := fs.String("ttl", "0", "TTL (e.g. 5s)")
		_ = fs.Parse(os.Args[2:])
		args := fs.Args()
		if len(args) < 2 {
			fmt.Println("usage: set <key> <value> [-ttl 5s]")
			return
		}
		key, val := args[0], args[1]
		owner := string(r.Owner(key))
		url := fmt.Sprintf("%s/cache/%s?ttl=%s", owner, key, *ttl)
		req, _ := http.NewRequest(http.MethodPut, url, strings.NewReader(val))
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)
		check(err)
		fmt.Println(resp.Status)

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("usage: get <key>")
			return
		}
		key := os.Args[2]
		owner := string(r.Owner(key))
		url := fmt.Sprintf("%s/cache/%s", owner, key)
		resp, err := http.Get(url)
		check(err)
		if resp.StatusCode != 200 {
			fmt.Println(resp.Status)
			return
		}
		b, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		fmt.Println(string(b))

	case "del":
		if len(os.Args) < 3 {
			fmt.Println("usage: del <key>")
			return
		}
		key := os.Args[2]
		owner := string(r.Owner(key))
		url := fmt.Sprintf("%s/cache/%s", owner, key)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)
		check(err)
		fmt.Println(resp.Status)

	default:
		usage()
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func usage() { fmt.Println("usage: set|get|del ... (env PEERS=comma_list)") }
func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "err:", err)
		os.Exit(1)
	}
}
