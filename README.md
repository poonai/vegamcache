# vegamcache
vegamcache is a distributed in-memory cache built using gossip protocol for golang.

# what is the difference between other distributed cache service?
In vegamcache, network calls are not used for retriving data from other nodes. Instead data will be replicated across the node using last win write in backgroud.

# seri why ?
Go is fun. I learned lot of thing regarding distributed system and also jobless.
# Drawback
- Can be used only in golang
- Consumes lot of main memory.(If you worring about memory, folks at google did a good job on group cache)

# Design

 
# Future Work
- shard the cache instead of storing in a single hashmap