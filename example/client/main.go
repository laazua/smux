package main

import (
	"flag"
	"fmt"

	"smux"
)

var (
	id      uint64
	message string
)

func main() {
	client := smux.NewClient("localhost:8886", &smux.JsonCode{})
	flag.Uint64Var(&id, "i", 0, "message id")
	flag.StringVar(&message, "m", "hello world", "send message")
	flag.Parse()

	msgData := &smux.Message{
		Id:   id,
		Body: []byte(message),
	}
	err := client.SendMessage(msgData)
	if err != nil {
		fmt.Printf("Send message error: %v\n", err)
		return
	}

	resp, err := client.RecvMessage()
	if err != nil {
		fmt.Printf("Recv message error: %v\n", err)
		return
	}
	fmt.Println("Recv message", resp.Id, string(resp.Body))

}
