package main

import (
	"bufio"
	"fmt"
	"os"
	chatclient "ueb10/client"
	"ueb10/protocol"
)

// Startet den Chat-Client und verbindet ihn mit dem Server.
// Nutzt die Serverdaten aus dem protocol-Paket und arbeitet mit einem ChatClient
// Liest den Benutzernamen über stdin ein, registriert ihn beim Server
// und startet danach die Lese- und Schreibschleife des Clients.
func main() {
	client, err := chatclient.NewChatClient(protocol.SERVER_HOST, protocol.SERVER_PORT)

	if err != nil {
		fmt.Printf(protocol.CLIENT_TEXT_CONNECTION_ERROR, err)
		return
	}

	defer client.Conn.Close()

	stdin := bufio.NewReader(os.Stdin)

	fmt.Println(protocol.CLIENT_TEXT_ASK_NAME)
	fmt.Println(protocol.CLIENT_TEXT_ALLOWED_NAME)

	registered := false

	for !registered {
		fmt.Print(protocol.CLIENT_TEXT_NAME_PROMPT)
		name, err := stdin.ReadString('\n')

		if err != nil {
			fmt.Println(protocol.CLIENT_TEXT_INPUT_ERROR)
			return
		}

		err = client.Register(name)

		if err != nil {
			fmt.Println("[Chat]", err)
		} else {
			fmt.Printf(protocol.CLIENT_TEXT_WELCOME, client.Name)
			registered = true
		}
	}

	fmt.Println(protocol.CLIENT_TEXT_COMMANDS)
	fmt.Println()

	done := make(chan bool, 1)

	// ReadLoop läuft parallel, damit Servernachrichten empfangen werden können,
	// während WriteLoop weiterhin Benutzereingaben liest.

	go client.ReadLoop(done)

	client.WriteLoop(stdin)
	<-done
}
