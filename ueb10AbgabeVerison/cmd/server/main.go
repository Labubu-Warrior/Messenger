package main

import (
	"fmt"
	"ueb10/protocol"
	chatserver "ueb10/server"
)

// Startet den Chat-Server.
// Erstellt einen neuen ChatServer und startet ihn auf dem Port aus dem protocol-Paket.
// Falls der Server nicht gestartet werden kann, wird der Fehler ausgegeben.
func main() {
	server := chatserver.NewChatServer()

	err := server.Start(protocol.SERVER_PORT)

	if err != nil {
		fmt.Printf("Serverfehler: %v\n", err)
	}
}
