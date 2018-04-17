/*
Copyright 2018 The vegamcache Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vegamcache

import (
	"bytes"
	"encoding/gob"
	"sync"
	"time"

	"github.com/weaveworks/mesh"
)

type Value struct {
	Data      []byte
	LastWrite int64
	Expiry    int64
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
	return c.mergeComplete(other.(*cache).copy().set)
}

func (c *cache) mergeComplete(other map[string]Value) mesh.GossipData {
	c.Lock()
	defer c.Unlock()
	for k, v := range other {
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

func (c *cache) get(key string) []byte {
	c.Lock()
	defer c.Unlock()
	if val, ok := c.set[key]; ok {
		if val.Expiry > time.Now().UnixNano() {
			return val.Data
		}
		return nil
	}
	return nil
}

func (c *cache) put(key string, val Value) {
	c.Lock()
	defer c.Unlock()
	c.set[key] = val
}
