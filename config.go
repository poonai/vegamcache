package vegamcache

import (
	"log"
	"math/rand"
	"os"

	"github.com/weaveworks/mesh"
)

type VegamConfig struct {
	Port     int
	Channel  string
	Password string
	NickName string
	Peers    []string
	PeerName string
	Host     string
	Logger   *log.Logger
}

func initConfig(config *VegamConfig) {
	if config.Port == 0 {
		config.Port = mesh.Port
	}
	if config.NickName == "" {
		name, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		config.NickName = name + string(rand.Int())
	}
	if config.Channel == "" {
		config.Channel = "default"
	}
	if config.PeerName == "" {
		config.PeerName = mustHardwareAddr()
	}
	if config.Host == "" {
		config.Host = "0.0.0.0"
	}
	if config.Logger == nil {
		config.Logger = log.New(os.Stdout, "", 0)
	}
}
