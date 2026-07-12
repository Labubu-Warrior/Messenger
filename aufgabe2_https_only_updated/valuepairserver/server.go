// Package valuepairserver implementiert den eigenen Wertepaar-Server.
package valuepairserver

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"aufgabe2/networkutil"
	"aufgabe2/plan"
	"aufgabe2/protocol"
	"aufgabe2/valuepair"
)

// Server berechnet freie Index-Wertepaare aus stets aktuellen XML-Daten.
type Server struct {
	PlanProvider *plan.Provider
}

// New erstellt einen Wertepaar-Server.
func New(provider *plan.Provider) *Server {
	return &Server{PlanProvider: provider}
}

// Serve verarbeitet beliebig viele voneinander unabhängige Verbindungen.
func (s *Server) Serve(listener net.Listener) error {
	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(connection)
	}
}

func (s *Server) handleConnection(connection net.Conn) {
	defer connection.Close()
	_ = connection.SetDeadline(time.Now().Add(protocol.VALUEPAIR_REQUEST_TIMEOUT))

	reader := bufio.NewReaderSize(connection, protocol.NETWORK_READ_BUFFER_BYTES)
	line, tooLong, readErr := networkutil.ReadDelimitedLine(
		reader,
		protocol.VALUEPAIR_END_MARKER,
		protocol.MAX_VALUEPAIR_LINE_BYTES,
	)
	if readErr != nil {
		writeError(connection, fmt.Errorf("Anfrage konnte nicht gelesen werden"))
		return
	}
	if tooLong {
		writeError(connection, fmt.Errorf("%s", protocol.MESSAGE_VALUEPAIR_REQUEST_TOO_LONG))
		return
	}

	ids, err := valuepair.ParseIDRequest(line)
	if err != nil {
		writeError(connection, err)
		return
	}

	currentPlan, err := s.PlanProvider.Load()
	if err != nil {
		writeError(connection, err)
		return
	}

	pairs, err := currentPlan.FreePairs(ids)
	if err != nil {
		writeError(connection, err)
		return
	}
	_, _ = fmt.Fprint(connection, valuepair.FormatPairs(pairs))
}

func writeError(connection net.Conn, err error) {
	message := strings.ReplaceAll(err.Error(), string(protocol.VALUEPAIR_END_MARKER), "")
	_, _ = fmt.Fprintf(connection, "%s%s%c", protocol.VALUEPAIR_ERROR_PREFIX, message, protocol.VALUEPAIR_END_MARKER)
}
