package vegamcache

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/weaveworks/mesh"
)

type Value struct {
	Data      string
	LastWrite int64
}
type cache struct {
	sync.Mutex
	set map[string]Value
}

var _ mesh.GossipData = &cache{}

func (c *cache) Encode() [][]byte {
	c.Lock()
	defer c.Unlock()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(c.set); err != nil {
		panic(err)
	}
	return [][]byte{buf.Bytes()}
}

func (c *cache) Merge(other mesh.GossipData) mesh.GossipData {
	c.Lock()
	defer c.Unlock()
	for k, v := range other.(*cache).set {
		val, ok := c.set[k]
		if ok && val.LastWrite < v.LastWrite {
			c.set[k] = v
			continue
		}
		c.set[k] = v
	}
	return c
}

func (c *cache) mergeDelta(set map[string]Value) (delta mesh.GossipData) {
	for k, v := range set {
		val, ok := c.set[k]
		if ok && val.LastWrite > v.LastWrite {
			delete(set, k)
			continue
		}
		c.set[k] = v
	}
	if len(set) == 0 {
		return nil
	}
	return &cache{set: set}
}

func (c *cache) mergeRecived(set map[string]Value) (recived mesh.GossipData) {
	for k, v := range set {
		val, ok := c.set[k]
		if ok && val.LastWrite > v.LastWrite {
			delete(set, k)
			continue
		}
		c.set[k] = v
	}
	if len(set) == 0 {
		return nil
	}
	return &cache{set: set}
}

func (c *cache) copy() *cache {
	return &cache{
		set: c.set,
	}
}

func (c *cache) get(key string) (val string) {
	c.Lock()
	defer c.Unlock()
	if val, ok := c.set[key]; ok {
		return val.Data
	}
	return ""
}

func (c *cache) put(key string, val Value) {
	c.Lock()
	defer c.Unlock()
	c.set[key] = val
}
