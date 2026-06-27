package server

import (
	"Merd_Chat/protocol"
	"bufio"
	"fmt"
	"net"

	"sort"
	"strings"
)

// ================= DATENSTRUKTUREN =================

// ConnectedClient repräsentiert einen angemeldeten Chat-Teilnehmer.
type ConnectedClient struct {
	Name string
	Conn net.Conn
	Send chan string // Gepufferter Kanal für ausgehende Nachrichten
}

// ChatServer orchestriert alle Clients asynchron über Channels (ohne Mutex!).
type ChatServer struct {
	clients   map[string]*ConnectedClient
	register  chan *registerRequest
	join      chan *ConnectedClient
	leave     chan *ConnectedClient
	broadcast chan string
	private   chan privateMessage
}

// Hilfsstrukturen für die internen Server-Ereignisse
type registerRequest struct {
	name     string
	conn     net.Conn
	response chan *ConnectedClient
}

type privateMessage struct {
	senderName string
	targetName string
	text       string
}

// ================= KONSTRUKTOR & START =================

// NewChatServer erstellt eine neue, leere Server-Instanz.
func NewChatServer() *ChatServer {
	return &ChatServer{
		clients:   make(map[string]*ConnectedClient),
		register:  make(chan *registerRequest),
		join:      make(chan *ConnectedClient),
		leave:     make(chan *ConnectedClient),
		broadcast: make(chan string),
		private:   make(chan privateMessage),
	}
}

// Start öffnet den TCP-Port und lauscht auf neue Verbindungen.
func (s *ChatServer) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("Server lauscht auf", address)

	// Starte die zentrale Event-Schleife im Hintergrund
	go s.run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Verbindungsfehler:", err)
			continue
		}
		go s.handleNewConnection(conn)
	}
}

// ================= DIE ZENTRALE EVENT-SCHLEIFE =================

// run verarbeitet alle Kanäle sequenziell. Hierdurch ist kein Mutex nötig!
func (s *ChatServer) run() {
	for {
		select {
		case req := <-s.register:
			// Prüfe ob Name eindeutig und gültig ist
			if _, exists := s.clients[req.name]; exists || req.name == "" || req.name == "System" {
				req.response <- nil
			} else {
				client := &ConnectedClient{
					Name: req.name,
					Conn: req.conn,
					Send: make(chan string, 100), // Puffer verhindert Blockieren
				}
				s.clients[req.name] = client
				req.response <- client
			}

		case client := <-s.join:
			s.sendListToClient(client)
			s.broadcastMessage(fmt.Sprintf(protocol.MSG_JOIN, client.Name), "System")

		case client := <-s.leave:
			if _, exists := s.clients[client.Name]; exists {
				delete(s.clients, client.Name)
				close(client.Send)
				client.Conn.Close()
				fmt.Printf("[%s] abgemeldet.\n", client.Name)
				s.broadcastMessage(fmt.Sprintf(protocol.MSG_LEAVE, client.Name), "System")
			}

		case msg := <-s.broadcast:
			for _, client := range s.clients {
				client.Send <- msg
			}

		case pm := <-s.private:
			if target, exists := s.clients[pm.targetName]; exists {
				target.Send <- fmt.Sprintf("[Privat von %s]: %s", pm.senderName, pm.text)
				// Auch dem Sender eine Bestätigung schicken
				if sender, ok := s.clients[pm.senderName]; ok {
					sender.Send <- fmt.Sprintf("[Privat an %s]: %s", pm.targetName, pm.text)
				}
			} else {
				if sender, ok := s.clients[pm.senderName]; ok {
					sender.Send <- fmt.Sprintf(protocol.ERR_USER_NOT_FOUND, pm.targetName)
				}
			}
		}
	}
}

// ================= CLIENT HANDLER =================

// handleNewConnection führt den Registrierungsprozess durch.
func (s *ChatServer) handleNewConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	fmt.Fprintf(conn, "%s%c", protocol.PROMPT_NAME, protocol.ENDSIGN)

	var client *ConnectedClient

	// Anmelde-Schleife (Eindeutigkeit prüfen)
	for client == nil {
		line, err := reader.ReadString(protocol.ENDSIGN)
		if err != nil {
			conn.Close()
			return
		}
		name := cleanString(line)

		// Anfrage an die zentrale Schleife schicken
		responseChan := make(chan *ConnectedClient)
		s.register <- &registerRequest{name: name, conn: conn, response: responseChan}
		client = <-responseChan

		if client == nil {
			fmt.Fprintf(conn, "%s%c", protocol.ERR_NAME_TAKEN, protocol.ENDSIGN)
		}
	}

	fmt.Printf("[%s] verbunden.\n", client.Name)
	client.Send <- fmt.Sprintf(protocol.MSG_WELCOME, client.Name)

	// Registrierung abgeschlossen, Client betritt den Chat
	s.join <- client

	// Start der Lese- und Schreib-Schleifen für diesen Client
	go s.writeLoop(client)
	s.readLoop(client, reader)
}

// readLoop empfängt kontinuierlich Nachrichten vom Client.
func (s *ChatServer) readLoop(client *ConnectedClient, reader *bufio.Reader) {
	// Sorgt für sicheren Logout bei Verbindungsabbruch
	defer func() { s.leave <- client }()

	for {
		line, err := reader.ReadString(protocol.ENDSIGN)
		if err != nil {
			return // Client hat die Verbindung getrennt
		}
		text := cleanString(line)
		if text == "" {
			continue
		}

		// Befehle verarbeiten
		parts := strings.SplitN(text, " ", 3)
		cmd := strings.ToLower(parts[0])

		switch cmd {
		case protocol.CMD_QUIT:
			return // Beendet die Schleife -> defer wird ausgeführt
		case protocol.CMD_LIST:
			s.sendListToClient(client)
		case protocol.CMD_PRIVATE:
			if len(parts) < 3 {
				client.Send <- protocol.ERR_PRIVATE_FORMAT
			} else {
				s.private <- privateMessage{senderName: client.Name, targetName: parts[1], text: parts[2]}
			}
		default:
			// Normale öffentliche Nachricht
			s.broadcast <- fmt.Sprintf("[%s]: %s", client.Name, text)
		}
	}
}

// writeLoop entleert den Send-Kanal und schreibt in den Socket.
func (s *ChatServer) writeLoop(client *ConnectedClient) {
	for msg := range client.Send {
		fmt.Fprintf(client.Conn, "%s%c", msg, protocol.ENDSIGN)
	}
}

// ================= HILFSFUNKTIONEN =================

func (s *ChatServer) sendListToClient(client *ConnectedClient) {
	var names []string
	for name := range s.clients {
		names = append(names, name)
	}
	sort.Strings(names)
	client.Send <- fmt.Sprintf(protocol.MSG_CLIENT_LIST, strings.Join(names, ", "))
}

func (s *ChatServer) broadcastMessage(msg string, sender string) {
	s.broadcast <- msg
}

func cleanString(s string) string {
	s = strings.TrimSuffix(s, string(protocol.ENDSIGN))
	return strings.TrimSpace(s)
}
