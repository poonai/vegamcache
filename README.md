# vegamcache
vegamcache is a distributed in-memory cache built using gossip protocol for golang.

# what is the difference between other distributed cache service?
In vegamcache, network calls are not used for retriving data from other nodes. Instead data will be replicated across the node using gossip in backgroud.

# seri why ?
Go is fun. I learned lot of thing regarding distributed system and also I'm jobless. Looking for internship. If anyone interested, do ping me at rbalajis25@gmail.com

# Drawback
- Can be used only in golang
- Consumes lot of main memory.(If you worring about memory, folks at google did a good job on group cache)
 
# Need to be done
- sharding the cache instead of storing it in a single hashmap
- adding expiry time
- small garbage collector to remove the expired value
- benchmarking against other cache service