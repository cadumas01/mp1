# mp1
A simple tcp multi-node network demonstration

# Usage:

- For each node listed in the config file, open an terminal instance and run: ``` go run main.go NODE_ID ``` where NODE_ID is the id that node as specified by config.txt

- Once all of the nodes are connected, send messages by entering: ```send DESTINATION_ID MESSAGE``` 


# Walkthrough

Each node is setup through the follwowing steps:
1. Querying the config file to get the min and max message send delays, and the ip address and port of the current node based on the ```id``` command line argument
2. Start the tcp server using ```net``` package's ```net.Listen(...)``` 
3. Create an empty map remote node id's (```string```) to incoming connections (```net.Conn```)
4. Start a goroutine to accept incoming dial requests from other nodes
    - Add these connections to the incoming connections map
    - For each connection, start a goroutine to handle each connection (reading anything that is writtten to that conneciton by the corresponding remote node)
5. Create a map of remote node id's to outgoing connections. Fill this map by dialing to each of the other nodes and adding each connection to the map
6. Once both the incoming and outgoing connection maps are complete (each node is conencted to all other nodes), each node can handle user input
    - Users can send messages from one node to a another node
        - Messages sent will have an artificial delay, randomly generated between the bounds specified by the config file.
    - Users can exit all nodes by typing ```$exit``` on one node

# Notes

# Assignement Directions
https://docs.google.com/document/d/1fB7Xudj6MQ690WjYsjvfp8CP_1eay3We959teEFEqIE/edit

# Resources
https://www.baeldung.com/cs/distributed-systems-guide
