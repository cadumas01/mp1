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
    - DO THIS

  - Node Dials to Another node once a message needs to be sent. Each node is both a server and client
    - Error checking if other node has not been initalize
    - Once Node A accepts Node B's dial, that means Node B has sent a message to node A so Node A should listen
      - This way don't need to keep a dictionary of all connections?

# Notes
## config
- IP address varies, port stays the same

# Assignement Directions
https://docs.google.com/document/d/1fB7Xudj6MQ690WjYsjvfp8CP_1eay3We959teEFEqIE/edit

# Resources
https://www.baeldung.com/cs/distributed-systems-guide
