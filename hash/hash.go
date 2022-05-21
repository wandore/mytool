package hash

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type Hash struct {
	circle        map[uint32]string
	virtualHashes []uint32
	clusterNodes  map[string]bool
	replicas      int
	hashFunc      func(data []byte) uint32
	mu            sync.RWMutex
}

func New(replicas int, hashFunc func(data []byte) uint32) *Hash {
	if hashFunc == nil {
		hashFunc = crc32.ChecksumIEEE
	}

	return &Hash{
		circle:        make(map[uint32]string, 0),
		virtualHashes: make([]uint32, 0),
		clusterNodes:  make(map[string]bool, 0),
		replicas:      replicas,
		hashFunc:      hashFunc,
		mu:            sync.RWMutex{},
	}
}

func (h *Hash) Add(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < h.replicas; i++ {
		virtualNode := node + ":" + strconv.Itoa(i)
		virtualNodeHash := h.getHash(virtualNode)
		h.virtualHashes = append(h.virtualHashes, virtualNodeHash)
		h.circle[virtualNodeHash] = node
	}

	h.sortVirtualNodes()

	h.clusterNodes[node] = true
}

func (h *Hash) Remove(node string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i := 0; i < h.replicas; i++ {
		virtualNode := node + ":" + strconv.Itoa(i)
		virtualNodeHash := h.getHash(virtualNode)

		index := sort.Search(len(h.virtualHashes), func(i int) bool {
			return h.virtualHashes[i] == virtualNodeHash
		})
		h.virtualHashes = append(h.virtualHashes[:index], h.virtualHashes[index+1:]...)

		delete(h.circle, virtualNodeHash)
	}

	h.sortVirtualNodes()

	delete(h.clusterNodes, node)
}

func (h *Hash) Match(key string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.circle) == 0 {
		return "", fmt.Errorf("hash circle is empty")
	}

	keyHash := h.getHash(key)
	index := sort.Search(len(h.virtualHashes), func(i int) bool {
		return h.virtualHashes[i] >= keyHash
	})
	if index >= len(h.virtualHashes) {
		index = 0
	}

	virtualHash := h.virtualHashes[index]

	return h.circle[virtualHash], nil
}

func (h *Hash) sortVirtualNodes() {
	sort.Slice(h.virtualHashes, func(i, j int) bool {
		return h.virtualHashes[i] < h.virtualHashes[j]
	})
}

func (h *Hash) getHash(data string) uint32 {
	return h.hashFunc([]byte(data))
}
