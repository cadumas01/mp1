package node

import (
	"bufio"
	"fmt"
	"mp1/configurations"
	"net"
	"os"
	"strconv"
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

	fmt.Println("After accept clinets")
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
		//skip first line that holds min and max delay

		lineArr := strings.Split(scanner.Text(), " ")
		if len(lineArr) == 3 {
			len := strconv.Itoa(len(lineArr))
			fmt.Println("Len of lineArr = " + len)
			fmt.Println("Here is lineArr: " + scanner.Text())
			if lineArr[0] != exclude_id {
				wg.Add(1)
				// ip + port

				// dial and add conn to map
				go connectTo(lineArr, connsMap, &wg)
			}
		}

	}

	wg.Wait()
	return connsMap
}

func connectTo(lineArr []string, conns_map map[string]net.Conn, wg *sync.WaitGroup) {
	address := lineArr[1] + ":" + lineArr[2]

	//Connect to port
	conn, err := net.Dial("tcp", address)

	for err != nil {
		fmt.Println("Dialing...")
		conn, err = net.Dial("tcp", address)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Broke out of loop, dial successful")
	// Delete
	if err != nil {
		panic(err)
	}

	// add to map
	conns_map[lineArr[0]] = conn

	fmt.Println("Client Successfully connected to  " + remoteConnIp(conn))
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

	//strip new line
	message = strings.Replace(message, "\n", "", -1)

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
	fmt.Println("Sent '" + message + "' to node" + destinationId + ", system time is" + time)

}

// HADNLING  Accepting connections and RECEIVING ///

func startServer(address string) (ln net.Listener) {
	ln, err := net.Listen("tcp", address)

	if err != nil {
		panic("Error listening")
	}

	fmt.Println("server started for " + ln.Addr().String())

	return ln
}

// Waits for client to connect and recieves message
func acceptClients(connections map[string]net.Conn, ln net.Listener) {

	fmt.Println("Inside acceptClients")
	// loop to allow function to accept all clients
	for {
		fmt.Println("Inside loop")
		// Waits for client to connect

		connChan := make(chan net.Conn)
		errChan := make(chan error)

		// Use of channels to return values from goroutine (ln.Accept())
		go func() {
			fmt.Println("About to accept")

			conn, err := ln.Accept() // NEED TO ADD goroutine????
			connChan <- conn
			errChan <- err
		}()

		// offloading channels
		conn := <-connChan
		err := <-errChan

		fmt.Println("Line 165")

		if err != nil {
			fmt.Println("About to panic")

			panic("error accepting")
		}

		acceptedIp := remoteConnIp(conn)
		acceptedId := configurations.QuerryConfig(acceptedIp, 1)[0]
		fmt.Println("acceptedIP : " + acceptedIp)
		// Find id from AcceptedIp - NOT DONE

		// Add connection to map
		connections[acceptedId] = conn //  CORRECT?, key must be id not ip
		fmt.Println("Just accepted " + acceptedIp + ", added id= " + acceptedId + " to connections map")

		go handleConnection(conn) // Bug here

	}

}

// fix bug
// Handles incoming messages for the node
func handleConnection(conn net.Conn) {

	fmt.Println("Inside handleConnection ")
	// loop to allow for many connection handling
	for {
		buf := make([]byte, bufSize)
		_, err := bufio.NewReader(conn).Read(buf)

		// if err is empty, we have a message and can print it
		if err == nil {
			message := string(buf)

			fmt.Println("Just read a message from " + remoteConnIp(conn))
			source := configurations.QuerryConfig(remoteConnIp(conn), 1)[1]
			unicastReceive(source, message)

		}
	}
	return
}

func unicastReceive(source string, message string) {
	time := time.Now().String()
	fmt.Println("Received '" + message + "' from node" + source + ", system time is" + time)
}

// Returns Conn's ip address
func remoteConnIp(conn net.Conn) string {
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		return addr.IP.String()
	} else {
		return ""
	}
}
