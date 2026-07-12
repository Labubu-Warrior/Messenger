// Package client enthält Netzwerkzugriff und grafische Oberfläche des Clients.
package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"aufgabe2/protocol"
	"aufgabe2/shared"
)

// Client hält eine dauerhafte TCP-Verbindung zum Hauptserver.
type Client struct {
	connection net.Conn
	reader     *bufio.Reader
	encoder    *json.Encoder
	mutex      sync.Mutex
	Codes      []string
}

// Connect verbindet sich und liest die automatisch gesendete Kürzelliste.
func Connect(address string) (*Client, error) {
	connection, err := net.DialTimeout(protocol.NETWORK_TYPE_TCP, address, protocol.CLIENT_CONNECT_TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("Hauptserver nicht erreichbar: %w", err)
	}

	client := &Client{
		connection: connection,
		reader:     bufio.NewReaderSize(connection, protocol.MAX_CLIENT_LINE_BYTES),
		encoder:    json.NewEncoder(connection),
	}
	response, err := client.readResponse()
	if err != nil {
		connection.Close()
		return nil, err
	}
	if response.Type == protocol.RESPONSE_ERROR {
		connection.Close()
		return nil, fmt.Errorf("%s", response.Message)
	}
	if response.Type != protocol.RESPONSE_WELCOME {
		connection.Close()
		return nil, fmt.Errorf("unerwartete Serverantwort")
	}
	client.Codes = append([]string(nil), response.Codes...)
	return client, nil
}

// RequestSuggestions sendet die vom Benutzer ausgewählten Kürzel.
func (c *Client) RequestSuggestions(codes []string) ([]shared.Suggestion, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	request := shared.ClientRequest{
		Type:  protocol.REQUEST_SUGGESTIONS,
		Codes: codes,
	}
	if err := c.encoder.Encode(request); err != nil {
		return nil, fmt.Errorf("Anfrage konnte nicht gesendet werden: %w", err)
	}

	response, err := c.readResponse()
	if err != nil {
		return nil, err
	}
	if response.Type == protocol.RESPONSE_ERROR {
		return nil, fmt.Errorf("%s", response.Message)
	}
	if response.Type != protocol.RESPONSE_SUGGESTIONS {
		return nil, fmt.Errorf("unerwartete Serverantwort")
	}
	return response.Suggestions, nil
}

func (c *Client) readResponse() (shared.ClientResponse, error) {
	var response shared.ClientResponse
	line, err := c.reader.ReadBytes(protocol.JSON_LINE_DELIMITER)
	if err != nil {
		return response, fmt.Errorf("Serverantwort konnte nicht gelesen werden: %w", err)
	}
	if err := json.Unmarshal(line, &response); err != nil {
		return response, fmt.Errorf("ungültige Serverantwort: %w", err)
	}
	return response, nil
}

// Close beendet die Verbindung sauber.
func (c *Client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_ = c.encoder.Encode(shared.ClientRequest{Type: protocol.REQUEST_QUIT})
	return c.connection.Close()
}
