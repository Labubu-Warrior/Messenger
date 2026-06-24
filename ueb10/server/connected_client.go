package server

import (
	"fmt"
	"net"
	"ueb10/protocol"
)

const SEND_BUFFER_SIZE = 16

// ConnectedClient repräsentiert einen am Server angemeldeten Chat-Teilnehmer.
type ConnectedClient struct {
	Name string
	Conn net.Conn
	Send chan string
}

// NewConnectedClient erstellt einen neuen serverseitigen Client.
func NewConnectedClient(name string, conn net.Conn) *ConnectedClient {
	return &ConnectedClient{
		Name: name,
		Conn: conn,
		Send: make(chan string, SEND_BUFFER_SIZE),
	}
}

// SendMessage legt eine Nachricht in den Sendekanal des Clients.
// Wenn der Sendepuffer voll ist, wird ein Fehler zurückgegeben.
func (c *ConnectedClient) SendMessage(message string) error {
	select {
	case c.Send <- message:
		return nil
	default:
		return fmt.Errorf(protocol.SERVER_ERROR_SEND_BUFFER_FULL, c.Name)
	}
}

// Close schließt die TCP-Verbindung des Clients.
func (c *ConnectedClient) Close() {
	c.Conn.Close()
}
