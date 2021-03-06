package main

import (
	"fmt"
	"mp1/configurations"
	"mp1/node"
	"os"
)

func usage() {
	fmt.Println("USAGE:\n\tTo start a node: go run main.go [NODE_ID]")
}

func main() {

	// Starting Node / CLI handling
	if len(os.Args) != 2 {
		usage()
		return
	}

	querry := configurations.QueryConfig(os.Args[1], 0)

	// if empty then id not found in config
	if len(querry) == 0 {
		fmt.Println("Machine ID not found in config file. Either add it to configs or initialize an ID already in config file")
		return
	}

	id := querry[0]
	node.StartNode(id)

	fmt.Println("Success")

}
