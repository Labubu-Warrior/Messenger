package server

import (
	"fmt"
	"net"
	"ueb10/protocol"
)

// ChatServer verwaltet alle angemeldeten Clients und verteilt Nachrichten.
// Die clients-Map wird nur in der Run-Schleife verändert.
// Dadurch braucht der Server keinen Mutex.
type ChatServer struct {
	register      chan RegisterRequest
	cancelPending chan string
	join          chan *ConnectedClient
	leave         chan *ConnectedClient
	broadcast     chan BroadcastMessage
	private       chan PrivateMessage
	listReq       chan *ConnectedClient

	clients      map[string]*ConnectedClient
	pendingNames map[string]struct{}
}

// NewChatServer erstellt einen neuen ChatServer mit allen benötigten Kanälen.
func NewChatServer() *ChatServer {
	return &ChatServer{
		register:      make(chan RegisterRequest),
		cancelPending: make(chan string),
		join:          make(chan *ConnectedClient),
		leave:         make(chan *ConnectedClient),
		broadcast:     make(chan BroadcastMessage),
		private:       make(chan PrivateMessage),
		listReq:       make(chan *ConnectedClient),
		clients:       make(map[string]*ConnectedClient),
		pendingNames:  make(map[string]struct{}),
	}
}

// Run ist die zentrale Server-Schleife.
// Sie verarbeitet alle Server-Ereignisse nacheinander.
func (s *ChatServer) Run() {
	for {
		select {
		case req := <-s.register:
			s.handleRegister(req)

		case name := <-s.cancelPending:
			s.handleCancelPending(name)

		case client := <-s.join:
			s.handleJoin(client)

		case client := <-s.leave:
			s.handleLeave(client)

		case message := <-s.broadcast:
			s.handleBroadcast(message)

		case pm := <-s.private:
			s.handlePrivateMessage(pm)

		case client := <-s.listReq:
			s.handleListRequest(client)
		}
	}
}

// Start öffnet den TCP-Port und nimmt neue Client-Verbindungen an.
func (s *ChatServer) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf(protocol.SERVER_TEXT_RUNNING, port)

	go s.Run()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf(protocol.SERVER_TEXT_ACCEPT_ERROR, err)
		} else {
			go s.HandleRegistration(conn)
		}
	}
}
