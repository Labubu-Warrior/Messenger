package main

import (
	"fmt"
	"ueb10/protocol"
	chatserver "ueb10/server"
)

func main() {
	server := chatserver.NewChatServer()

	err := server.Start(protocol.SERVER_PORT)

	if err != nil {
		fmt.Printf("Serverfehler: %v\n", err)
	}
}
