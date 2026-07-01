package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"ueb10/protocol"
)

// Methode von ChatServer
// Übernimmt die Anmeldung eines neu verbundenen Clients.
// Erwartet eine offene TCP-Verbindung und liest daraus den gewünschten Namen.
// Registriert den Namen über den zentralen Serverkanal und startet den Client erst,
// wenn der Server die Registrierung bestätigt hat.
// Bei ungültigem Namen, belegtem Namen oder Verbindungsfehler wird passend reagiert.
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

		if isLineTooLong(line) {
			fmt.Fprintf(conn, "%s\n", protocol.RESPONSE_NAME_INVALID)
			continue
		}

		result := s.tryRegister(conn, line)

		if result.OK {
			client = result.Client
			_, err = fmt.Fprintf(conn, "%s\n", protocol.RESPONSE_NAME_OK)

			if err != nil {
				s.cancelPending <- client.Name
				conn.Close()
				return
			}

			registered = true
		} else {
			fmt.Fprintf(conn, "%s\n", result.Reason)
		}
	}

	s.startClient(client, reader)
}

// Methode von ChatServer
// Sendet einen Registrierungsversuch an den zentralen Serverloop.
// Erwartet Verbindung und gewünschten Namen.
// Liefert das Ergebnis der Registrierung über einen eigenen Antwortkanal zurück.
// Dadurch entscheidet nur der Serverloop, ob ein Name frei ist.
func (s *ChatServer) tryRegister(conn net.Conn, name string) RegisterResponse {
	respChan := make(chan RegisterResponse, 1)

	s.register <- RegisterRequest{
		Name: name,
		Conn: conn,
		Resp: respChan,
	}

	return <-respChan
}

// Methode von ChatServer
// Nimmt einen erfolgreich registrierten Client in den Chat auf.
// Meldet den Client über den join-Kanal beim Serverloop an.
// Startet danach getrennte Schleifen für eingehende und ausgehende Nachrichten.
func (s *ChatServer) startClient(client *ConnectedClient, reader *bufio.Reader) {
	s.join <- client

	go s.WriteLoop(client)
	go s.ReadLoop(client, reader)
}

// Methode von ChatServer
// Liest dauerhaft Nachrichten eines angemeldeten Clients.
// Erwartet den Client und dessen Reader für die TCP-Verbindung.
// Beendet den Client bei Verbindungsabbruch oder QUIT.
// Gültige Eingaben werden zur weiteren Verarbeitung weitergegeben.
func (s *ChatServer) ReadLoop(client *ConnectedClient, reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			s.leave <- client
			return
		}

		if isLineTooLong(line) {
			s.sendOrLog(client, protocol.ERROR_LINE_TOO_LONG)
			continue
		}

		if s.isQuitCommand(line) {
			s.sendOrLog(client, protocol.MESSAGE_GOODBYE)
			s.leave <- client
			return
		}

		s.handleClientInput(client, line)
	}
}

// Prüft, ob eine empfangene Zeile die erlaubte Maximallänge überschreitet.
// Entfernt vorher nur die Zeilenendzeichen, damit Enter nicht mitgezählt wird.
// Liefert true, wenn die Nutzdaten länger als erlaubt sind.
func isLineTooLong(line string) bool {
	return len(strings.TrimRight(line, "\r\n")) > protocol.MAX_LINE_LENGTH
}

// Methode von ChatServer
// Prüft, ob die empfangene Zeile ein QUIT-Kommando enthält.
// Erwartet eine rohe Eingabezeile vom Client.
// Liefert true, wenn das erste Wort dem Protokollbefehl QUIT entspricht.
func (s *ChatServer) isQuitCommand(line string) bool {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)

	if len(parts) == 0 {
		return false
	}

	command := strings.ToUpper(parts[0])
	return command == protocol.COMMAND_QUIT
}

// Methode von ChatServer
// Wertet eine bereinigte Client-Eingabe als Serverkommando aus.
// Erwartet den sendenden Client und die empfangene Zeile.
// Leitet MSG, PRIVMSG und LIST an die passende Verarbeitung weiter.
// Unbekannte Befehle werden dem Client als Fehler gemeldet.
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

// Methode von ChatServer
// Verarbeitet ein öffentliches Chatkommando.
// Erwartet den sendenden Client und die bereits zerlegte Eingabe.
// Sendet gültige Nachrichten über den broadcast-Kanal an den Serverloop.
// Fehlt der Nachrichtentext, erhält nur der Sender eine Fehlermeldung.
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

// Methode von ChatServer
// Verarbeitet ein privates Chatkommando.
// Erwartet den sendenden Client, den Empfänger und den Nachrichtentext in parts.
// Sendet gültige private Nachrichten über den private-Kanal an den Serverloop.
// Bei fehlenden Angaben erhält nur der Sender einen Nutzungshinweis.
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

// Methode von ChatServer
// Sendet ausgehende Nachrichten an einen bestimmten Client.(client.Conn)
// Erwartet einen angemeldeten Client mit Send-Kanal und TCP-Verbindung.
// Liest Nachrichten aus client.Send und schreibt sie in die Verbindung.
// Bei Sendefehler oder geschlossenem Send-Kanal wird die Verbindung geschlossen.

func (s *ChatServer) WriteLoop(client *ConnectedClient) {
	for msg := range client.Send {
		_, err := fmt.Fprintf(client.Conn, "%s\n", msg) // пишет msg в TCP-соединение конкретного клиента.

		if err != nil {
			//
			client.Close()
			return
		}
	}

	client.Close()
}
