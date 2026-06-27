package client

import (
	"Merd_Chat/protocol"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// ChatClient repräsentiert die Nutzerseite der Verbindung.
type ChatClient struct {
	conn net.Conn
}

// Connect baut die Verbindung zum Server auf.
func Connect(host, port string) (*ChatClient, error) {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return nil, err
	}
	return &ChatClient{conn: conn}, nil
}

// Run startet den Client, empfängt den Login-Prompt und beginnt den Chat.
func (c *ChatClient) Run() {
	defer c.conn.Close()

	serverReader := bufio.NewReader(c.conn)
	terminalReader := bufio.NewReader(os.Stdin)

	// Starte einen Hintergrund-Thread für eingehende Server-Nachrichten
	go func() {
		for {
			msg, err := serverReader.ReadString(protocol.ENDSIGN)
			if err != nil {
				fmt.Println("\n[Verbindung zum Server verloren]")
				os.Exit(0)
			}
			msg = strings.TrimSuffix(msg, string(protocol.ENDSIGN))
			fmt.Println(msg)
		}
	}()

	// Haupt-Thread: Warte auf Tastatureingaben und sende sie an den Server
	for {
		text, _ := terminalReader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text != "" {
			fmt.Fprintf(c.conn, "%s%c", text, protocol.ENDSIGN)
		}

		if strings.ToLower(text) == protocol.CMD_QUIT {
			break
		}
	}
}
