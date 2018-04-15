package vegamcache

import (
	"github.com/weaveworks/mesh"
)

type peer struct {
	cache *cache
}

var _ mesh.Gossiper = &peer{}

func (p *peer) Gossip() (complete mesh.GossipData) {
	return p.cache.Copy()
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
