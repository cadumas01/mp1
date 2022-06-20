package node

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"mp1/configurations"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var bufSize int = 2048
var minDelay int
var maxDelay int
var exitCode string = "$exit"
var CONFIG string = "config.txt"

// Starts a node
// address = ip + port
func StartNode(id string) {

	// Delay boundaries
	minDelay, maxDelay = configurations.GetDelayBounds()

	// Initalize Listener
	lineArr := configurations.QuerryConfig(id, 0)
	address := lineArr[1] + ":" + lineArr[2]

	fmt.Println("address is " + address)
	ln := startServer(address)

	// connctions maps ip?? to connection
	in_conns := make(map[string]net.Conn)
	fmt.Println("len of map = " + strconv.Itoa(len(in_conns)))

	// Listen for connection
	// Accept connections, get their address with conn.RemoteAddr() and add to dictionary of connections

	var wgAccept sync.WaitGroup

	// total lines - 1 (delays) - 1 (self) = total number of connections
	wgAccept.Add(countLines(CONFIG) - 2)

	go acceptClients(in_conns, ln, &wgAccept)

	// Try to Dial into other listeners
	outConnsMap := OutConnsMap(id)

	wgAccept.Wait()
	// Wait for input and send messages
	fmt.Println("All nodes connected. Send a message with: send [DESTINATION ID] [MESSAGE]")

	for {
		destination, message := handleInput(outConnsMap)
		if destination != "" {
			go unicastSend(outConnsMap, destination, message)
		}
	}
}

// Handling Outgoing connection to Other Nodes ///

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

	if err != nil {
		panic(err)
	}

	// add to map
	conns_map[lineArr[0]] = conn

	fmt.Println("Client Successfully connected to  " + remoteConnIp(conn))
	wg.Done()
}

// First two words are commands, all other words are part of message
// Handles exit codes
func handleInput(connsMap map[string]net.Conn) (destination string, message string) {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	//strip new line
	text = strings.Replace(text, "\n", "", -1)

	// exit all nodes
	if text == "$exit" {
		for _, conn := range connsMap {
			_, err := conn.Write([]byte("$exit"))
			if err != nil {
				panic(err)
			}
		}

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

func unicastSend(connsMap map[string]net.Conn, destinationId string, message string) {
	conn := connsMap[destinationId]

	if conn == nil {
		fmt.Println("Invalid destination id, try again")
		return
	}

	// Send notification
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("Sent '" + message + "' to node " + destinationId + ", system time is " + now)

	// Artificial Delay before actually writing to channel
	time.Sleep(time.Duration(getDelay()) * time.Millisecond)

	// Write Message over tcp channel
	_, err := conn.Write([]byte(message))

	if err != nil {
		panic("Error writing message")
	}
}

// HADNLING  Accepting connections and RECEIVING ///

// Starts server for other nodes to connect to
func startServer(address string) (ln net.Listener) {
	ln, err := net.Listen("tcp", address)

	if err != nil {
		panic("Error listening")
	}

	fmt.Println("server started for " + ln.Addr().String())

	return ln
}

// Waits for client to connect and recieves message
func acceptClients(connections map[string]net.Conn, ln net.Listener, wgAccept *sync.WaitGroup) {
	// loop to allow function to accept all clients
	for {
		connChan := make(chan net.Conn)
		errChan := make(chan error)

		// Use of channels to return values from goroutine (ln.Accept())
		// Not sure if this is correct, test without it
		go func() {
			fmt.Println("About to accept")

			conn, err := ln.Accept() // NEED TO ADD goroutine????
			fmt.Println("Newly accepted conn is : " + remoteConnIp(conn))
			connChan <- conn
			errChan <- err
		}()

		// offloading channels
		conn := <-connChan
		err := <-errChan

		if err != nil {
			panic("error accepting")
		}

		acceptedIp := remoteConnIp(conn)
		fmt.Println("acceptedIp: " + acceptedIp)
		fmt.Println("Remote addr: " + conn.RemoteAddr().String())
		fmt.Println("IDK: " + conn.RemoteAddr().(*net.TCPAddr).AddrPort().String())
		fmt.Println("Local addr: " + conn.LocalAddr().String())

		acceptedId := configurations.QuerryConfig(acceptedIp, 1)[0]

		// Add connection to map
		connections[acceptedId] = conn
		fmt.Println("Just accepted " + acceptedIp + ", added id= " + acceptedId + " to connections map")

		go handleConnection(conn)
		wgAccept.Done()
	}
}

// Handles incoming messages for the node
func handleConnection(conn net.Conn) {
	// loop to allow for many connection handling
	for {
		buf := make([]byte, bufSize)
		_, err := bufio.NewReader(conn).Read(buf)

		// if err is empty, we have a message and can print it
		if err == nil {

			message := string(bytes.Trim(buf, "\x00")) //trims buf of empty bytes

			message = strings.TrimSpace(message) // maybe delete?
			source := configurations.QuerryConfig(remoteConnIp(conn), 1)[0]

			unicastReceive(source, message)
		}
	}
}

func unicastReceive(source string, message string) {
	if message == exitCode {
		fmt.Println("Exiting")
		os.Exit(0)
	}

	time := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("Received '" + message + "' from node " + source + ", system time is: " + time)
}

// Util///

// Returns Conn's ip address
func remoteConnIp(conn net.Conn) string {
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		return addr.IP.String()
	} else {
		return ""
	}
}

func getDelay() int {
	diff := maxDelay - minDelay
	return minDelay + rand.Intn(diff)
}

func countLines(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	return lineCount
}
