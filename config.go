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
