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
	"encoding/json"
	"net/http"
)

type UpdateRequest struct {
	Peers []string `json:"peers"`
}

type UpdateResponse struct {
	Updated          bool   `json:"updated"`
	ErrorDescription string `json:"error_description"`
}

func marshal(data interface{}) []byte {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return marshaledData
}
func updateHandler(v *vegam) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var peers UpdateRequest
		defer r.Body.Close()
		if r.Method != "PATCH" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(marshal(UpdateResponse{
				Updated:          false,
				ErrorDescription: "http method not allowed",
			}))
			return
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&peers)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(marshal(UpdateResponse{
				Updated:          false,
				ErrorDescription: err.Error(),
			}))
			return
		}
		v.router.ConnectionMaker.InitiateConnections(peers.Peers, true)
		w.Write(marshal(UpdateResponse{
			Updated: true,
		}))
	}

}
func ListenServer(v *vegam, port string) error {
	http.HandleFunc("/update", updateHandler(v))
	return http.ListenAndServe(port, nil)
}
