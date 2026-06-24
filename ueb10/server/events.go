package server

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"ueb10/protocol"
)

const VALID_NAME_PATTERN = `^[A-Za-z0-9_-]+$`

var validNamePattern = regexp.MustCompile(VALID_NAME_PATTERN)

func (s *ChatServer) handleRegister(req RegisterRequest) {
	name := strings.TrimSpace(req.Name)
	_, exists := s.clients[name]

	if !isValidClientName(name) {
		req.Resp <- RegisterResponse{
			OK:     false,
			Reason: protocol.RESPONSE_NAME_INVALID,
			Client: nil,
		}
	} else if exists {
		req.Resp <- RegisterResponse{
			OK:     false,
			Reason: protocol.RESPONSE_NAME_TAKEN,
			Client: nil,
		}
	} else {
		client := NewConnectedClient(name, req.Conn)
		s.clients[name] = client

		req.Resp <- RegisterResponse{
			OK:     true,
			Reason: protocol.RESPONSE_NAME_OK,
			Client: client,
		}
	}
}

func isValidClientName(name string) bool {
	return validNamePattern.MatchString(name)
}

func (s *ChatServer) handleJoin(client *ConnectedClient) {
	names := s.getClientNames(client.Name)

	s.sendOrLog(client, protocol.RESPONSE_NAME_OK)
	s.sendOrLog(client, protocol.MESSAGE_CLIENTS_PREFIX+strings.Join(names, ","))

	systemText := fmt.Sprintf(protocol.TEXT_SYSTEM_JOIN, client.Name)
	s.sendToAllExcept(client.Name, systemText)

	fmt.Printf(protocol.SERVER_TEXT_CLIENT_CONNECTED, client.Name)
}

func (s *ChatServer) handleLeave(client *ConnectedClient) {
	_, exists := s.clients[client.Name]

	if exists {
		delete(s.clients, client.Name)
		close(client.Send)

		systemText := fmt.Sprintf(protocol.TEXT_SYSTEM_LEAVE, client.Name)
		s.sendToAllExcept(client.Name, systemText)

		fmt.Printf(protocol.SERVER_TEXT_CLIENT_DISCONNECTED, client.Name)
	}
}

func (s *ChatServer) handleBroadcast(message BroadcastMessage) {
	out := fmt.Sprintf(protocol.TEXT_PUBLIC_MSG, message.From, message.Text)
	s.sendToAllExcept(message.From, out)
}

func (s *ChatServer) handlePrivateMessage(pm PrivateMessage) {
	receiver, exists := s.clients[pm.To]

	if exists {
		out := fmt.Sprintf(protocol.TEXT_PRIVATE_MSG, pm.From, pm.Text)
		s.sendOrLog(receiver, out)
	} else {
		sender, ok := s.clients[pm.From]

		if ok {
			s.sendOrLog(sender, protocol.ERROR_RECEIVER_NOT_FOUND)
		}
	}
}

func (s *ChatServer) handleListRequest(client *ConnectedClient) {
	names := s.getClientNames("")
	s.sendOrLog(client, protocol.MESSAGE_CLIENTS_PREFIX+strings.Join(names, ","))
}

func (s *ChatServer) getClientNames(excludeName string) []string {
	names := make([]string, 0, len(s.clients))

	for name := range s.clients {
		if excludeName == "" || name != excludeName {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names
}

func (s *ChatServer) sendToAllExcept(excludeName string, message string) {
	for name, client := range s.clients {
		if name != excludeName {
			s.sendOrLog(client, message)
		}
	}
}

func (s *ChatServer) sendOrLog(client *ConnectedClient, message string) {
	err := client.SendMessage(message)

	if err != nil {
		fmt.Printf(protocol.SERVER_TEXT_SEND_FAILED, client.Name, err)
	}
}
