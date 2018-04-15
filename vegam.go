package vegamcache

import (
	"net"
	"time"

	"github.com/weaveworks/mesh"
)

type vegam struct {
	gossip mesh.Gossip
	peer   *peer
	router *mesh.Router
}

func NewVegam(vc *VegamConfig) *vegam {
	initConfig(vc)
	peername, err := mesh.PeerNameFromString(vc.PeerName)
	if err != nil {
		panic(err)
	}
	router, err := mesh.NewRouter(
		mesh.Config{
			Port:               vc.Port,
			ProtocolMinVersion: mesh.ProtocolMinVersion,
			Password:           []byte(vc.Password),
			Host:               vc.Host,
			PeerDiscovery:      true,
			TrustedSubnets:     []*net.IPNet{},
		},
		peername,
		vc.NickName,
		mesh.NullOverlay{},
		vc.Logger,
	)
	if err != nil {
		panic(err)
	}
	peer := &peer{
		cache: &cache{
			set: make(map[string]LastWrite),
		},
	}
	gossip, err := router.NewGossip(vc.Channel, peer)
	if err != nil {
		panic(err)
	}
	router.Start()
	router.ConnectionMaker.InitiateConnections(vc.Peers, true)
	return &vegam{
		gossip: gossip,
		peer:   peer,
		router: router,
	}
}

func (v *vegam) Get(key string) (val string) {
	val = v.peer.cache.Get(key)
	return
}

func (v *vegam) Put(key, val string) {
	lw := time.Now().Unix()
	v.peer.cache.Put(key, val, lw)
	tempCache := &cache{
		set: make(map[string]LastWrite),
	}
	tempCache.set[key] = LastWrite{
		Value:     val,
		LastWrite: lw,
	}
	v.gossip.GossipBroadcast(tempCache)
}

func (v *vegam) Stop() {
	v.router.Stop()
}
