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
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := &cache{
		set: make(map[string]Value),
	}
	c.put("foo", Value{
		Data:      "bar",
		Expiry:    time.Now().Add(5 * time.Second).UnixNano(),
		LastWrite: time.Now().UnixNano(),
	})
	val, _ := c.get("foo")
	if val.(string) != "bar" {
		t.Fatalf("expected bar got %s", val)
	}
}

func TestExpiry(t *testing.T) {
	c := &cache{
		set: make(map[string]Value),
	}
	c.put("foo", Value{
		Data:      "bar",
		Expiry:    time.Now().Add(5 * time.Second).UnixNano(),
		LastWrite: time.Now().UnixNano(),
	})
	time.Sleep(5 * time.Second)
	_, exist := c.get("foo")
	if exist {
		t.Fatalf("value has to be exipred after 5 seconds")
	}
}

func TestMerges(t *testing.T) {
	c := &cache{
		set: make(map[string]Value),
	}
	c.put("foo", Value{
		Data:      "bar",
		Expiry:    time.Now().Add(5 * time.Second).UnixNano(),
		LastWrite: time.Now().UnixNano(),
	})
	completedData := c.mergeComplete(map[string]Value{"bar": Value{Data: "foo"}})
	if len(completedData.(*cache).set) != 2 {
		t.Fatalf("merge complete has to give the entire data")
	}
	for _, val := range []string{"foo", "bar"} {
		_, ok := completedData.(*cache).set[val]
		if !ok {
			t.Fatalf("%s not exist in set", val)
		}
	}
	recivedData := c.mergeRecived(map[string]Value{"bar": Value{Data: "foo"}})
	if len(recivedData.(*cache).set) != 1 {
		t.Fatalf("merge recived has to reture only recived data")
	}
	_, ok := recivedData.(*cache).set["bar"]
	if !ok {
		t.Fatal("bar not exist in set")
	}
}
