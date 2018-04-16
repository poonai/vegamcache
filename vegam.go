package vegamcache

import (
	"net"
	"time"

	"github.com/weaveworks/mesh"
)

type vegam struct {
	gossip  mesh.Gossip
	peer    *peer
	router  *mesh.Router
	actions chan<- func()
	peers   []string
	stop    chan int
}

func NewVegam(vc *VegamConfig) (*vegam, error) {
	initConfig(vc)
	peername, err := mesh.PeerNameFromString(vc.PeerName)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	peer := &peer{
		cache: &cache{
			set: make(map[string]Value),
		},
	}
	gossip, err := router.NewGossip(vc.Channel, peer)
	if err != nil {
		return nil, err
	}
	return &vegam{
		gossip: gossip,
		peer:   peer,
		router: router,
		peers:  vc.Peers,
		stop:   make(chan int),
	}, nil
}

func (v *vegam) Start() {
	actions := make(chan func())
	v.actions = actions
	v.router.Start()
	v.router.ConnectionMaker.InitiateConnections(v.peers, true)
	go v.loop(actions)
}

func (v *vegam) loop(actions <-chan func()) {
	for {
		select {
		case f := <-actions:
			f()
		case <-v.stop:
			return
		}
	}
}

func (v *vegam) Stop() {
	v.stop <- 1
	v.router.Stop()
}

func (v *vegam) Get(key string) (val []byte) {
	val = v.peer.cache.get(key)
	return
}

func (v *vegam) Put(key string, val []byte) {
	tempVal := Value{
		Data:      val,
		LastWrite: time.Now().Unix(),
	}
	v.peer.cache.put(key, tempVal)
	tempCache := &cache{
		set: make(map[string]Value),
	}
	tempCache.set[key] = tempVal
	v.actions <- func() {
		v.gossip.GossipBroadcast(tempCache)
	}
}
