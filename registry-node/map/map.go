package regmap

import (
	"fmt"
	"math/rand"

	"github.com/brendenehlers/go-distributed-cache/registry-node"
)

var (
	ErrMapSize     = fmt.Errorf("no registry nodes available")
	ErrNodeInvalid = fmt.Errorf("invalid node")
)

// key: url
type RegistryMap struct {
	values map[string]*registry.RegistryEntry
}

func New() *RegistryMap {
	return &RegistryMap{
		values: make(map[string]*registry.RegistryEntry),
	}
}

func (r *RegistryMap) Register(url string) error {
	r.values[url] = &registry.RegistryEntry{
		Url: url,
	}

	return nil
}

func (r *RegistryMap) Unregister(url string) error {
	delete(r.values, url)
	return nil
}

func (r *RegistryMap) GetNode() (*registry.RegistryEntry, error) {
	if len(r.values) == 0 {
		return nil, ErrMapSize
	}

	node := r.getRandomNode()

	if node == nil {
		return nil, ErrNodeInvalid
	}

	return node, nil
}

func (r *RegistryMap) getRandomNode() *registry.RegistryEntry {
	if len(r.values) == 0 {
		return nil
	}

	index := 0
	rIndex := rand.Intn(len(r.values))
	var node *registry.RegistryEntry

	for _, val := range r.values {
		if index == rIndex {
			node = val
			break
		}
		index++
	}

	return node
}
