package server

import (
	"fmt"
	"net"
	"ueb10/protocol"
)

// Maximale Anzahl von Nachrichten, die gleichzeitig,
// im Send-Kanal eines Clients zwischengespeichert werden können.
const SEND_BUFFER_SIZE = 16

// сама очередь сообщений этого клиента
// ConnectedClient repräsentiert einen am Server angemeldeten Chat-Teilnehmer.
type ConnectedClient struct {
	Name string
	Conn net.Conn
	Send chan string
}

// Erstellt einen neuen ConnectedClient für einen registrierten Benutzer.
// Erwartet den Benutzernamen und eine geöffnete TCP-Verbindung.
// Initialisiert zusätzlich den Sendekanal mit einer festen Puffergröße.
// Liefert den vollständig initialisierten ConnectedClient zurück.
func NewConnectedClient(name string, conn net.Conn) *ConnectedClient {
	return &ConnectedClient{
		Name: name,
		Conn: conn,
		Send: make(chan string, SEND_BUFFER_SIZE),
	}
}

// Methode von ConnectedClient
// Übergibt eine Nachricht an den Sendekanal(Send) des Clients.
// Erwartet den Nachrichtentext als Parameter.
// Liefert nil bei erfolgreichem Einfügen in den Sendepuffer
// oder einen Fehler, falls der Sendepuffer bereits voll ist.
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
