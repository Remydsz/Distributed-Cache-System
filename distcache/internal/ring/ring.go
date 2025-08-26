package ring

import (
	"crypto/sha1"
	"sort"
	"strconv"
	"sync"
)

type NodeID string

type vnode struct {
	hash uint32
	id   NodeID
}

type Ring struct {
	mu sync.RWMutex
	vn []vnode // sorted by hash
}

// New builds a ring with "replicas" vnodes per node.
func New(nodes []NodeID, replicas int) *Ring {
	r := &Ring{}
	r.Rebuild(nodes, replicas)
	return r
}

func (r *Ring) Rebuild(nodes []NodeID, replicas int) {
	var vs []vnode
	for _, n := range nodes {
		for rep := 0; rep < replicas; rep++ {
			h := hash32(string(n) + "#" + strconv.Itoa(rep))
			vs = append(vs, vnode{hash: h, id: n})
		}
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].hash < vs[j].hash })
	r.mu.Lock()
	r.vn = vs
	r.mu.Unlock()
}

// Owner returns the node responsible for a key.
func (r *Ring) Owner(key string) NodeID {
	h := hash32(key)
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.vn) == 0 {
		return ""
	}
	i := sort.Search(len(r.vn), func(i int) bool { return r.vn[i].hash >= h })
	if i == len(r.vn) {
		i = 0
	}
	return r.vn[i].id
}

func hash32(s string) uint32 {
	sum := sha1.Sum([]byte(s))
	return uint32(sum[0])<<24 | uint32(sum[1])<<16 | uint32(sum[2])<<8 | uint32(sum[3])
}
