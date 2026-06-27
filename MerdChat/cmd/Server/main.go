package main

import (
	server "Merd_Chat/ServerLogik"
	"Merd_Chat/protocol"
	"fmt"
)

func main() {
	chatServer := server.NewChatServer()
	address := protocol.SERVER_HOST + ":" + protocol.SERVER_PORT

	err := chatServer.Start(address)
	if err != nil {
		fmt.Println("Fehler beim Starten des Servers:", err)
	}
}
