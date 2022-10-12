package main

import (
	"fmt"
	"net"
	"os"
)

const NB_PART = 2

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	writeToFile("log_coordinator.txt", "begin_commit")

	// Listen for two participant
	participants := make([]net.Conn, 0)

	for i := 0; i < NB_PART; i++ {
		participant, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		participants = append(participants, participant)
	}

	fmt.Println(participants)

	// send prepare message
	prepare_message := Message{Action: "PREPARE", Payload: ""}
	prepare_participants(participants, prepare_message)

	// handle prepare response
	handleParticipantsResponses(participants)

	getAcks(participants)

	writeToFile("log_coordinator.txt", "end_of_transaction")
}

func prepare_participants(participants []net.Conn, message Message) {
	for _, participant := range participants {
		send(participant, message)
	}
}

func handleParticipantsResponses(participants []net.Conn) {
	is_a_participant_abort := false

	for _, participant := range participants {
		prepare_response := receive(participant)

		if prepare_response.Action == "ABORT" {
			is_a_participant_abort = true
		}
	}

	for _, participant := range participants {
		if is_a_participant_abort {
			handleAbort(participant)
		} else {
			handleCommit(participant)
		}
	}

	if is_a_participant_abort {
		writeToFile("log_coordinator.txt", "abort")
	} else {
		writeToFile("log_coordinator.txt", "commit")
	}
}

func handleAbort(conn net.Conn) {
	// send the global abort to other participants
	message := Message{Action: "ABORT", Payload: ""}
	send(conn, message)
}

func handleCommit(conn net.Conn) {
	// send the global commit to other participants
	message := Message{Action: "COMMIT", Payload: "chat"}
	send(conn, message)
}

func getAcks(participants []net.Conn) {
	for _, participant := range participants {
		receive(participant)
		participant.Close()
	}
}
