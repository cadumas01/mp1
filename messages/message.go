package messages

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type messageString struct {
	To      string `json:"To"`
	From    string `json:"From"`
	Date    string `json:"Date"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

type Message struct {
	To      string
	From    string
	Date    time.Time
	Title   string
	Content string
}

func newMessage(mS messageString) *Message {
	const longForm = "January 2, 2006 3:04pm (MST)"
	date, err := time.Parse(longForm, mS.Date)

	if err != nil {
		fmt.Println("Invalid time")
	}

	return &Message{mS.To, mS.From, date, mS.Title, mS.Content}
}

func fromCLI() *Message {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("MESSAGE GUIDELINES: End each field by hitting enter")

	args := []string{"To: ", "From: ", "Title: ", "Content: "}
	var inputs [4]string
	for i := 0; i < len(args); i++ {
		fmt.Print(args[i])
		input, _ := reader.ReadString('\n')

		// trim newline
		input = strings.TrimSuffix(input, "\n")
		inputs[i] = input
	}

	return &Message{inputs[0], inputs[1], time.Now(), inputs[2], inputs[3]}
}

func fromJson(path string) *Message {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	// Must unmarshall the json object
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var mS messageString
	json.Unmarshal(byteValue, &mS)

	jsonFile.Close()

	return newMessage(mS)
}

func ConstructMessage() *Message {
	reader := bufio.NewReader(os.Stdin)
	var m *Message
	fmt.Println("Construct a message from CLI? ('Y' for CLI, else to use message1.json)")
	input, _ := reader.ReadString('\n')

	if input == "Y\n" {
		m = fromCLI()
	} else {
		m = fromJson("message1.json")
	}
	return m
}

func (m Message) String() string {
	return fmt.Sprintf("\nTo: %s\nFrom: %s\nDate: %s\nTitle: %s\nContent: %s\n\n",
		m.To, m.From, m.Date.String(), m.Title, m.Content)
}
