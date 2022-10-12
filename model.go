package main

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3334"
	CONN_TYPE = "tcp"
)

type Message struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
}
