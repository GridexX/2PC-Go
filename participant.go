package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
)

func main() {
	// Listen for incoming connections.
	//Connect to the coordinator
	conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error connection:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.

	// handle prepare
	// prepare_message := receive(conn)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// Receive a prepare message
	handleRequest(conn)

	// Receive a commit or abort message
	handleRequest(conn)
	send(conn, Message{Action: "ACK", Payload: ""})
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	message := receive(conn)

	switch message.Action {
	case "ABORT":
		handleAbort()
	case "COMMIT":
		handleCommit(message.Payload)
	case "PREPARE":
		handlePrepare(conn)
	}
}

func getNameFile() string {
	//Write the commit in log file
	name_log := "1"
	name := os.Getenv("NUM")
	if name != "" {
		name_log = name
	}
	return "log_participant_" + name_log + ".txt"
}

func handleAbort() {
	//Write the abort in log file
	writeToFile(getNameFile(), "abort")
}

func handleCommit(payload string) {

	writeToFile(getNameFile(), "commit")

	// writeToFile("log_participant.txt", "commit")
	//Write the new entry in a new file
	writeToFile("data.txt", "Insert :"+payload)
}

func handlePrepare(conn net.Conn) {

	writeToFile(getNameFile(), "ready")
	//random number between 0 and 1
	random_number := rand.Float64()
	action := "COMMIT"
	if random_number < 0.01 {
		action = "ABORT"
	}
	send(conn, Message{Action: action, Payload: ""})
}
