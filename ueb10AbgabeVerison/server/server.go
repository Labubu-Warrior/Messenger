package server

import (
	"fmt"
	"net"
	"ueb10/protocol"
)

// Verwaltet den Zustand des Chat-Servers.
// Die Kanäle bilden die Schnittstelle zwischen Client-Goroutinen und Run-Schleife.
// clients enthält aktive Clients, pendingNames reserviert Namen während der Anmeldung.
// Da nur Run diese Maps verändert, ist kein Mutex nötig.
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

// Erstellt einen leeren ChatServer.
// Initialisiert alle Kanäle für Serverereignisse sowie die Maps für aktive
// und noch nicht vollständig angemeldete Clients.
// Liefert den vorbereiteten Server zurück.
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

// Verarbeitet alle Serverereignisse zentral und nacheinander.
// Empfängt Registrierungen, Joins, Leaves, Nachrichten und Listenanfragen
// über die Serverkanäle und ruft die passenden Handler auf.
// Diese Schleife besitzt die Server-Maps und synchronisiert damit den Zustand.
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

// Startet den TCP-Server auf dem angegebenen Port.
// Öffnet den Listener, startet die zentrale Run-Schleife
// und nimmt danach dauerhaft neue Client-Verbindungen an.
// Liefert einen Fehler zurück, wenn der Port nicht geöffnet werden kann.
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
