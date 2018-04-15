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
