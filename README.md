# vegamcache 
vegamcache is a distributed in-memory cache built using gossip protocol for golang.

# what is the difference between other distributed cache service?
In vegamcache, network calls are not used for retriving data for each Get. Instead data will be replicated across the node using gossip in backgroud.

Expired keys are removed on gossip instead of having a seperate GC.
# seri why ?
Go is fun. I learned lot of thing regarding distributed system and also I'm jobless. Looking for internship. If anyone interested, do ping me at rbalajis25@gmail.com

# Drawback
- Can be used only in golang
- Consumes lot of main memory.(If you worring about memory, folks at google did a good job on group cache)
 
# Need to be done
- sharding the cache instead of storing it in a single hashmap
- benchmarking against other cache service

# Example

## Clustered Cache

```go
vg, err := vegamcache.NewVegam(&vegamcache.VegamConfig{Port: 8087,
            PeerName: "00:00:00:00:00:01",
            Peers: []string{"remoteip1:port","remoteip2:port"},
			Logger:   log.New(ioutil.Discard, "", 0)})
vg.Start()
defer vg.Stop()
if err != nil {
    panic(err)
}
vg.Put("foo", "bar", time.Second*200)
val, exist := vg.Get("foo")
if exist {
    fmt.Println(val)
}
```

## Single Node Cache

```go
vg := vegamcache.NewCache()
vg.Put("foo", "bar", time.Second*200)
val, exist := vg.Get("foo")
if exist {
    fmt.Println(val)
}
```

# Contribution

Feel free to send PR. :)