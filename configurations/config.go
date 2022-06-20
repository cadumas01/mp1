package configurations

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const CONFIG string = "config.txt"

// Queries config file to return arguments for a given line based on ID
// Returns [ID, IP, PORT]
// keyType may be: id=0, ip=1, port=2
func QueryConfig(key string, keyType int) []string {
	file, err := os.Open(CONFIG)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var empty []string
	for scanner.Scan() {

		lineArr := strings.Split(scanner.Text(), " ")
		if lineArr[keyType] == key {
			return lineArr
		}
	}

	return empty
}

func GetDelayBounds() (min int, max int) {
	file, err := os.Open(CONFIG)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	lineArr := strings.Split(scanner.Text(), " ")

	if len(lineArr) != 2 {
		panic("Invalid delay boundaries")
	}

	min, err = strconv.Atoi(lineArr[0])
	if err != nil {
		panic(err)
	}

	max, err = strconv.Atoi(lineArr[1])
	if err != nil {
		panic(err)
	}

	return
}
