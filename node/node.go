package node

import (
	"bufio"
	"fmt"
	"mp1/configurations"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var bufSize int = 2048

// address = ip + port
func StartNode(id string) {
	// Initalize Listener
	lineArr := configurations.QuerryConfig(id, 0)
	address := lineArr[1] + ":" + lineArr[2]

	fmt.Println("address is " + address)
	ln := startServer(address)

	// connctions maps ip?? to connection
	in_conns := make(map[string]net.Conn)

	// Listen for connection
	// Accept connections, get their address with conn.RemoteAddr() and add to dictionary of connections
	go acceptClients(in_conns, ln)

	// Try to Dial into other listeners
	outConnsMap := OutConnsMap(id)

	// Wait for input and send messages
	for {
		destination, message := handleInput()
		if destination != "" {
			unicastSend(outConnsMap, destination, message)
		}
	}
}

// HANDLING Connecting to Other Nodes ///
// Creates map of outgoing connections = {id1: conn1, ...} without exclude_id (self)
// Each conn is null to start
func OutConnsMap(exclude_id string) map[string]net.Conn {
	// Do not return until all conncetions have been formed

	var wg sync.WaitGroup

	file, err := os.Open("config.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	connsMap := make(map[string]net.Conn)

	for scanner.Scan() {

		lineArr := strings.Split(scanner.Text(), " ")
		if lineArr[0] != exclude_id {
			wg.Add(1)
			// ip + port

			// dial and add conn to map
			go connect(lineArr, connsMap, &wg)
		}
	}

	wg.Wait()
	return connsMap
}

func connect(lineArr []string, conns_map map[string]net.Conn, wg *sync.WaitGroup) {
	address := lineArr[1] + ":" + lineArr[2]

	//Connect to port
	conn, err := net.Dial("tcp", address)

	if err != nil {
		panic(err)
	}

	// add to map
	conns_map[lineArr[0]] = conn

	fmt.Println("Client Successfully connected from " + conn.LocalAddr().String())
	wg.Done()
}

// First two words are commands, all other words are part of message
func handleInput() (destination string, message string) {
	// handle input
	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	if text == "exit" {
		fmt.Println("Exiting")
		os.Exit(0)
	}

	textArr := strings.Split(text, " ")

	if len(textArr) < 3 {
		fmt.Println("Invalid command")
		return
	}

	destination = textArr[1]
	message = strings.Join(textArr[2:], " ")
	return
}

// FINISH
func unicastSend(connsMap map[string]net.Conn, destinationId string, message string) {

	// Dealing with Message construction
	conn := connsMap[destinationId]

	// Send Message
	fmt.Println("Sending a message...")

	// Write Message over tcp channel
	_, err := conn.Write([]byte(message))

	if err != nil {
		panic("Error writing message")
	}

	//Sent “Hello” to process 2, system time is ­­­­­­­­­­­­­XXX
	time := time.Now().String()
	fmt.Println("Sent'" + message + "' to node" + destinationId + ", system time is" + time)

}

// HADNLING  Accepting connections and RECEIVING ///

func startServer(address string) (ln net.Listener) {
	ln, err := net.Listen("tcp", address)

	if err != nil {
		panic("Error listening")
	}

	fmt.Println("server started for " + address)

	return ln
}

// Waits for client to connect and recieves message
func acceptClients(connections map[string]net.Conn, ln net.Listener) {

	fmt.Println("Inside acceptClients")
	// loop to allow function to accept all clients
	for {
		// Waits for client to connect
		conn, err := ln.Accept() // NEED TO ADD goroutine????
		fmt.Println("Line 165")

		if err != nil {
			panic("error accepting")
		}

		acceptedIp := conn.RemoteAddr().String()
		fmt.Println("acceptedIP : " + acceptedIp)
		// Find id from AcceptedIp - NOT DONE

		// Add connection to map
		connections[acceptedIp] = conn // NOT CORRECT, key must be id not ip
		fmt.Println("Just accepted " + acceptedIp)

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {

	fmt.Println("INside handleConnection")
	// loop to allow for many connection handling
	for {
		buf := make([]byte, bufSize)
		_, err := bufio.NewReader(conn).Read(buf)

		if err != nil {
			panic(err)
		}

		message := string(buf)

		source := configurations.QuerryConfig(conn.RemoteAddr().String(), 1)[1]
		unicastReceive(source, message)

		// Send confiramtion message
		//conn.Write([]byte(Confirmation))
		//fmt.Println("Confirmation sent, server exiting...")

	}
	return
}

func unicastReceive(source string, message string) {
	time := time.Now().String()
	fmt.Println("Received'" + message + "' from node" + source + ", system time is" + time)
}
