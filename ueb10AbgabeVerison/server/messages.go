package server

import "net"

// PrivateMessage beschreibt eine Nachricht an genau einen Empfänger.
type PrivateMessage struct {
	From string
	To   string
	Text string
}

// BroadcastMessage beschreibt eine öffentliche Nachricht an alle anderen Clients.
type BroadcastMessage struct {
	From string
	Text string
}

// RegisterRequest wird an die Server-Schleife gesendet,
// wenn ein Client einen Namen anmelden möchte.
type RegisterRequest struct {
	Name string
	Conn net.Conn
	Resp chan RegisterResponse
}

// RegisterResponse ist die Antwort des Servers auf einen Registrierungsversuch.
type RegisterResponse struct {
	OK     bool
	Reason string
	Client *ConnectedClient
}
