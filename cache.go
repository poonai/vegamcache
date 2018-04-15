package vegamcache

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/weaveworks/mesh"
)

type LastWrite struct {
	Value     string
	LastWrite int64
}
type cache struct {
	sync.Mutex
	set map[string]LastWrite
}

// state implements GossipData.
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

func (c *cache) mergeDelta(set map[string]LastWrite) (delta mesh.GossipData) {
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

func (c *cache) mergeRecived(set map[string]LastWrite) (recived mesh.GossipData) {
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

func (c *cache) Copy() *cache {
	return &cache{
		set: c.set,
	}
}

func (c *cache) Get(key string) (val string) {
	c.Lock()
	defer c.Unlock()
	if val, ok := c.set[key]; ok {
		return val.Value
	}
	return ""
}

func (c *cache) Put(key, val string, writetime int64) {
	c.Lock()
	defer c.Unlock()
	c.set[key] = LastWrite{
		Value:     val,
		LastWrite: writetime,
	}
}
