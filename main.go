package main

import (
	"fmt"
	"mp1/configs"
	"os"
)

func usage() {
	fmt.Println("USAGE:\n\tTo start a node: go run main.go [NODE_ID]\n\tSend a message with: send [DESTINATION] [MESSAGE]")
}

func main() {

	// Starting Node / CLI handling
	if len(os.Args) != 2 {
		usage()
		return
	}

	querry := configs.QuerryConfig(os.Args[1], 0)
	//querry := configs.QuerryConfig(os.Args[1])

	// if empty then id not found in config
	if len(querry) == 0 {
		fmt.Println("Machine ID not found in config file. Either add it to configs or initialize an ID already in config file")
		return
	}

	fmt.Println("Success")

}
