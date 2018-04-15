package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/sch00lb0y/vegamcache"
)

func main() {
	var get = flag.String("action", "true", "define usage")
	flag.Parse()
	if *get == "true" {
		fmt.Print("GET IS RUNNING")
		vg := vegamcache.NewVegam(&vegamcache.VegamConfig{Port: 8087,
			PeerName: "00:00:00:00:00:01"})
		// ticker := time.NewTimer(5 * time.Second)
		// go func() {
		// 	for _ = range ticker.C {
		// 		fmt.Print("Getting")
		// 		val := vg.Get("foo")
		// 		if val == "bar" {
		// 			fmt.Print(val)
		// 			ticker.Stop()
		// 		} else {
		// 			fmt.Print("sup")
		// 		}
		// 	}
		// }()

		ticker := time.NewTicker(500 * time.Millisecond)
		go func() {
			for _ = range ticker.C {
				fmt.Println("Getting")
				val := vg.Get("foo")
				if val == "bar" {
					fmt.Print(val)
					ticker.Stop()
				} else {
					fmt.Println(val)
				}
			}
		}()

		// Tickers can be stopped like timers. Once a ticker
		// is stopped it won't receive any more values on its
		// channel. We'll stop ours after 1600ms.
		stop := make(chan int)
		<-stop
	} else {

		fmt.Print("PUT IS RUNNING")
		vg := vegamcache.NewVegam(&vegamcache.VegamConfig{Port: 8081,
			PeerName: "00:00:00:00:00:02",
			Peers:    []string{"localhost:8087"}})
		vg.Put("foo", "bar")
		fmt.Print(vg.Get("foo"))
		stop := make(chan int)
		<-stop
	}

}
