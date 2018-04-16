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
	"github.com/weaveworks/mesh"
)

type peer struct {
	cache *cache
}

var _ mesh.Gossiper = &peer{}

func (p *peer) Gossip() (complete mesh.GossipData) {
	return p.cache.copy()
}

func (p *peer) OnGossip(buf []byte) (delta mesh.GossipData, err error) {
	set, err := decodeSet(buf)
	if err != nil {
		return nil, err
	}
	return p.cache.mergeDelta(set), nil
}

func (p *peer) OnGossipBroadcast(src mesh.PeerName, buf []byte) (recived mesh.GossipData, err error) {
	set, err := decodeSet(buf)
	if err != nil {
		return nil, err
	}
	return p.cache.mergeRecived(set), nil
}

func (p *peer) OnGossipUnicast(src mesh.PeerName, buf []byte) error {
	set, err := decodeSet(buf)
	if err != nil {
		return err
	}
	p.cache.mergeRecived(set)
	return nil
}
