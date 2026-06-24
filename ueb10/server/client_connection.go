package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"ueb10/protocol"
)

// HandleRegistration liest Namenseingaben eines neuen Clients,
// bis der Server einen gültigen Namen akzeptiert.
func (s *ChatServer) HandleRegistration(conn net.Conn) {
	reader := bufio.NewReader(conn)
	registered := false
	var client *ConnectedClient

	for !registered {
		line, err := reader.ReadString('\n')

		if err != nil {
			conn.Close()
			return
		}

		result := s.tryRegister(conn, line)

		if result.OK {
			client = result.Client
			registered = true
		} else {
			fmt.Fprintf(conn, "%s\n", result.Reason)
		}
	}

	s.startClient(client, reader)
}

func (s *ChatServer) tryRegister(conn net.Conn, name string) RegisterResponse {
	respChan := make(chan RegisterResponse, 1)

	s.register <- RegisterRequest{
		Name: name,
		Conn: conn,
		Resp: respChan,
	}

	return <-respChan
}

func (s *ChatServer) startClient(client *ConnectedClient, reader *bufio.Reader) {
	go s.WriteLoop(client)
	go s.ReadLoop(client, reader)

	s.join <- client
}

// ReadLoop liest dauerhaft Befehle von einem angemeldeten Client.
func (s *ChatServer) ReadLoop(client *ConnectedClient, reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			s.leave <- client
			return
		}

		if s.isQuitCommand(line) {
			s.sendOrLog(client, protocol.MESSAGE_GOODBYE)
			s.leave <- client
			return
		}

		s.handleClientInput(client, line)
	}
}

func (s *ChatServer) isQuitCommand(line string) bool {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)

	if len(parts) == 0 {
		return false
	}

	command := strings.ToUpper(parts[0])
	return command == protocol.COMMAND_QUIT
}

func (s *ChatServer) handleClientInput(client *ConnectedClient, line string) {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)

	if len(parts) == 0 {
		return
	}

	command := strings.ToUpper(parts[0])

	switch command {
	case protocol.COMMAND_MESSAGE:
		s.handleMessageCommand(client, parts)

	case protocol.COMMAND_PRIVATE:
		s.handlePrivateCommand(client, parts)

	case protocol.COMMAND_LIST:
		s.listReq <- client

	default:
		s.sendOrLog(client, protocol.ERROR_UNKNOWN_COMMAND)
	}
}

func (s *ChatServer) handleMessageCommand(client *ConnectedClient, parts []string) {
	if len(parts) < protocol.MIN_MESSAGE_PARTS {
		s.sendOrLog(client, protocol.ERROR_MESSAGE_MISSING)
	} else {
		text := strings.Join(parts[1:], " ")

		s.broadcast <- BroadcastMessage{
			From: client.Name,
			Text: text,
		}
	}
}

func (s *ChatServer) handlePrivateCommand(client *ConnectedClient, parts []string) {
	if len(parts) < protocol.MIN_PRIVATE_PARTS {
		s.sendOrLog(client, protocol.ERROR_PRIVATE_USAGE)
	} else {
		to := parts[1]
		text := strings.Join(parts[2:], " ")

		s.private <- PrivateMessage{
			From: client.Name,
			To:   to,
			Text: text,
		}
	}
}

// WriteLoop schreibt alle Servernachrichten an einen einzelnen Client.
func (s *ChatServer) WriteLoop(client *ConnectedClient) {
	for msg := range client.Send {
		_, err := fmt.Fprintf(client.Conn, "%s\n", msg)

		if err != nil {
			client.Close()
			return
		}
	}

	client.Close()
}
