// Package server implementiert den Hauptserver für GUI-Clients.
package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"aufgabe2/networkutil"
	"aufgabe2/protocol"
	"aufgabe2/shared"
)

// Server verwaltet unabhängige Client-Verbindungen.
type Server struct {
	Service *TerminService
}

// New erstellt einen Hauptserver.
func New(service *TerminService) *Server {
	return &Server{Service: service}
}

// Serve akzeptiert Clients und verarbeitet jeden Client in einer eigenen Goroutine.
func (s *Server) Serve(listener net.Listener) error {
	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleClient(connection)
	}
}

func (s *Server) handleClient(connection net.Conn) {
	defer connection.Close()
	encoder := json.NewEncoder(connection)

	codes, err := s.Service.Codes()
	if err != nil {
		s.writeResponse(encoder, shared.ClientResponse{
			Type:    protocol.RESPONSE_ERROR,
			Message: err.Error(),
		})
		return
	}
	if !s.writeResponse(encoder, shared.ClientResponse{
		Type:  protocol.RESPONSE_WELCOME,
		Codes: codes,
	}) {
		return
	}

	reader := bufio.NewReaderSize(connection, protocol.NETWORK_READ_BUFFER_BYTES)
	for {
		_ = connection.SetReadDeadline(time.Now().Add(protocol.CLIENT_READ_TIMEOUT))

		line, tooLong, readErr := networkutil.ReadDelimitedLine(
			reader,
			protocol.JSON_LINE_DELIMITER,
			protocol.MAX_CLIENT_LINE_BYTES,
		)
		if readErr != nil {
			return
		}
		if tooLong {
			if !s.writeResponse(encoder, shared.ClientResponse{
				Type:    protocol.RESPONSE_ERROR,
				Message: protocol.MESSAGE_REQUEST_TOO_LONG,
			}) {
				return
			}
			continue
		}

		var request shared.ClientRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			if !s.writeResponse(encoder, shared.ClientResponse{
				Type:    protocol.RESPONSE_ERROR,
				Message: protocol.MESSAGE_INVALID_REQUEST_FORMAT,
			}) {
				return
			}
			continue
		}

		switch strings.ToUpper(strings.TrimSpace(request.Type)) {
		case protocol.REQUEST_SUGGESTIONS:
			suggestions, err := s.Service.Suggestions(request.Codes)
			if err != nil {
				if !s.writeResponse(encoder, shared.ClientResponse{
					Type:    protocol.RESPONSE_ERROR,
					Message: err.Error(),
				}) {
					return
				}
				continue
			}
			if !s.writeResponse(encoder, shared.ClientResponse{
				Type:        protocol.RESPONSE_SUGGESTIONS,
				Suggestions: suggestions,
			}) {
				return
			}

		case protocol.REQUEST_QUIT:
			s.writeResponse(encoder, shared.ClientResponse{
				Type:    protocol.RESPONSE_GOODBYE,
				Message: protocol.MESSAGE_CONNECTION_CLOSED,
			})
			return

		default:
			if !s.writeResponse(encoder, shared.ClientResponse{
				Type:    protocol.RESPONSE_ERROR,
				Message: protocol.MESSAGE_UNKNOWN_COMMAND,
			}) {
				return
			}
		}
	}
}

func (s *Server) writeResponse(encoder *json.Encoder, response shared.ClientResponse) bool {
	if err := encoder.Encode(response); err != nil {
		fmt.Printf("Antwort konnte nicht gesendet werden: %v\n", err)
		return false
	}
	return true
}
