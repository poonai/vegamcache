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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/balajijinnah/vegamcache"
)

func main() {
	var get = flag.String("action", "true", "get or put")
	flag.Parse()
	if *get == "true" {
		fmt.Print("GET IS RUNNING")
		vg, _ := vegamcache.NewVegam(&vegamcache.VegamConfig{Port: 8087,
			PeerName: "00:00:00:00:00:01",
			Logger:   log.New(ioutil.Discard, "", 0)})
		vg.Start()
		ticker := time.NewTicker(500 * time.Millisecond)
		go func() {
			for _ = range ticker.C {
				val, exist := vg.Get("foo")
				if exist && val.(string) == "bar" {
					fmt.Println(val)
					fmt.Println("value got from another node")
					ticker.Stop()
				} else {
					fmt.Println(val)
				}
			}
		}()
		stop := make(chan int)
		<-stop
	} else {
		fmt.Print("PUT IS RUNNING")
		vg, _ := vegamcache.NewVegam(&vegamcache.VegamConfig{Port: 8082,
			PeerName: "00:00:00:00:00:02"})
		vg.Start()
		vg.Put("foo", "bar", time.Second*200)
		fmt.Print(vg.Get("foo"))
		go vegamcache.ListenServer(vg, ":8000")
		req, _ := http.NewRequest("PATCH", "http://localhost:8000/update", bytes.NewBuffer(
			[]byte(`{
				"peers":[
					"localhost:8087"
				]
			}`),
		))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		stop := make(chan int)
		<-stop
	}
	// vg := vegamcache.NewCache()
	// vg.Put("foo", "bar", 5*time.Second)
	// val, _ := vg.Get("foo")
	// fmt.Print(val)
}
