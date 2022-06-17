package node

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var bufSize int = 2048

// address = ip + port
func Start_Node(address string) {
	// Initalize Listener
	ln := startServer(address)

	// connctions maps ip?? to connection
	// connections := make(map[string]net.Conn)

	// Listen for connection
	// Accept connections, get their address with conn.RemoteAddr() and add to dictionary of connections
	go acceptClient(ln)

	// Try to Dial into other listeners
}

func startServer(address string) (ln net.Listener) {
	ln, err := net.Listen("tcp", address)

	if err != nil {
		panic("Error listening")
	}

	fmt.Println("server started for " + address)

	return ln
}

// Waits for client to connect and recieves message
func acceptClient(ln net.Listener) {

	// Waits for client to connect
	conn, err := ln.Accept()

	if err != nil {
		panic("error accepting")
	}

	acceptedIp := conn.RemoteAddr().String()
	fmt.Println("Just accepted " + acceptedIp)

	handleConnection(conn)
	//connections[acceptedIp] = conn
}

func handleConnection(conn net.Conn) {

	buf := make([]byte, bufSize)
	_, err := bufio.NewReader(conn).Read(buf)

	if err != nil {
		panic(err)
	}

	message := string(buf)

	source = mp
	unicast_receive()

	// Send confiramtion message
	//conn.Write([]byte(Confirmation))
	//fmt.Println("Confirmation sent, server exiting...")

	return
}

func unicast_receive(source string, message string) {
	time := time.Now().String()
	fmt.Println("Received'" + message + "' from" + source + ", system time is" + time)
}
