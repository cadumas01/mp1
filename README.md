# mp1
A simple network demonstration


# Ideas: 
- 1 server that each network application connects to?
  - This could be its own proces?
  - Main process that sets everything up

- Each node/process communicates directly with other nodes
  - Trying this first
  - Each node is named by its ID and created with IP address and port number specfied by the config file
  - Are all nodes initialized at the start? Or is they only connected upon request
    - Could try to dial to all other nodes immedaitely (block on goroutines) until the other nodes are up


# Resources
https://www.baeldung.com/cs/distributed-systems-guide
