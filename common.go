package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

func send(conn net.Conn, message Message) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err.Error())
		os.Exit(1)
	}
	conn.Write(jsonMessage)
	log.Default().Println("Message sent.", message)
}

func receive(conn net.Conn) Message {
	// bufio.NewReader(conn).ReadString('\n')

	buf := make([]byte, 2048)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	if n > 0 {
		log.Default().Println("Message received.", string(buf))
	}
	var message Message
	err = json.Unmarshal(buf[:n], &message)

	if err != nil {
		fmt.Println("Error reading JSON:", err.Error())
	}

	return message
}

func writeToFile(fileName, content string) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(content + "\n")
	if err != nil {
		log.Fatal(err)
	}
}
