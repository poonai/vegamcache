package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/sch00lb0y/vegamcache"
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
				val := vg.Get("foo")
				if val == "bar" {
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
		vg, _ := vegamcache.NewVegam(&vegamcache.VegamConfig{Port: 8081,
			PeerName: "00:00:00:00:00:02",
			Peers:    []string{"localhost:8087"}})
		vg.Start()
		vg.Put("foo", "bar")
		fmt.Print(vg.Get("foo"))
		stop := make(chan int)
		<-stop
	}
}
