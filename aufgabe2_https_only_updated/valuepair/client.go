package valuepair

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"aufgabe2/protocol"
	"aufgabe2/shared"
)

// Client stellt Anfragen an den eigenen Wertepaar-Server.
type Client struct {
	Address string
}

// NewClient erstellt einen Wertepaar-Client.
func NewClient(address string) *Client {
	return &Client{Address: address}
}

// RequestPairs sendet Mitarbeiter-IDs und liest freie Index-Wertepaare.
func (c *Client) RequestPairs(ids []int) ([]shared.ValuePair, error) {
	request, err := FormatIDRequest(ids)
	if err != nil {
		return nil, err
	}

	connection, err := net.DialTimeout(protocol.NETWORK_TYPE_TCP, c.Address, protocol.VALUEPAIR_REQUEST_TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("Wertepaar-Server nicht erreichbar: %w", err)
	}
	defer connection.Close()
	_ = connection.SetDeadline(time.Now().Add(protocol.VALUEPAIR_REQUEST_TIMEOUT))

	if _, err := fmt.Fprint(connection, request); err != nil {
		return nil, fmt.Errorf("Wertepaar-Anfrage konnte nicht gesendet werden: %w", err)
	}

	reader := bufio.NewReaderSize(connection, protocol.MAX_VALUEPAIR_LINE_BYTES)
	response, err := reader.ReadString(protocol.VALUEPAIR_END_MARKER)
	if err != nil {
		return nil, fmt.Errorf("Wertepaar-Antwort konnte nicht gelesen werden: %w", err)
	}
	return ParsePairs(response)
}
