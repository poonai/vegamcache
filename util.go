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
	"net"
)

func mustHardwareAddr() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		if s := iface.HardwareAddr.String(); s != "" {
			return s
		}
	}
	panic("no valid network interfaces")
}

func decodeSet(buf []byte) (map[string]Value, error) {
	var set map[string]Value
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&set); err != nil {
		return nil, err
	}
	return set, nil
}
