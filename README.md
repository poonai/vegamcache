# vegamcache
vegamcache is a distributed in-memory cache built using gossip protocol for golang.

# what is the difference between other distributed cache service?
In vegamcache, network calls are not used for retriving data from other nodes. Instead data will be replicated across the node using last win write in backgroud.

# Drawback
- Can be used only in golang
- Consumes lot of main memory.(If you worring about memory, folks at google did a good job on group cache)

# Design

 